package sprinter

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// ControlGroup defines a cgroups v2 namespace and related resource limits
type ControlGroup struct {
	Name        string
	Limits      []Limit
	SysProcAttr *syscall.SysProcAttr
}

// Limit is a single cgroup v2 parameter and value
type Limit struct {
	Var   string
	Value string
}

type RunArgs struct {
	Dir              string               // working dir
	WrapperPath      string               // path to cmd wrapper executable
	SPA              *syscall.SysProcAttr // any SysProcAttr values we want applied to the running process
	Buffer           *SafeBuffer          // where should the output go?
	Cleanup          func()               // called when the process ends so we can clean up resources
	ControlGroupPath string               // path to cgroup
	Cmd              string               // command to execute
	Args             []string             // command args
}

// Driver breaks out the cgroup functionality so the tests don't need elevated privileges.
type Driver interface {
	MakeCgroup(dir string, cg ControlGroup) error
	Run(ra RunArgs) (*exec.Cmd, *StatusTracker, <-chan bool, error)
}

// Linux driver handles cgroup operations on a linux system
type LinuxDriver struct{}

func (l *LinuxDriver) MakeCgroup(dir string, cg ControlGroup) error {
	// makes cgroups in the v2 style of /sys/fs/cgroup/<ppid>
	if err := os.MkdirAll(dir, 0644); err != nil {
		if errors.Is(err, fs.ErrExist) {
			log.Println("cgroup", dir, "already exists")
		} else {
			log.Println("failed to create cgroup for", dir, ":", err)
			return err
		}
	}

	// Now we can write the actual cgroup limits ...
	for _, l := range cg.Limits {
		lf := filepath.Join(dir, l.Var)
		if err := os.WriteFile(lf, []byte(l.Value), fs.FileMode(os.O_APPEND)); err != nil {
			return err
		}
	}
	return nil
}

func (l *LinuxDriver) Run(ra RunArgs) (*exec.Cmd, *StatusTracker, <-chan bool, error) {
	args := append([]string{ra.ControlGroupPath, ra.Cmd}, ra.Args...)
	log.Println("running process:", ra.WrapperPath, args)
	cmd := exec.Command(ra.WrapperPath, args...)
	cmd.Stdout = ra.Buffer
	cmd.Stderr = ra.Buffer
	cmd.Dir = ra.Dir
	// NOTE: experimental - we may need to do this in the wrapper.
	cmd.SysProcAttr = ra.SPA // linux namespace controls

	done := make(chan bool)
	tracker := &StatusTracker{}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error running the exec.Command - %s\n", err)
		ra.Cleanup()
		return nil, nil, nil, err
	}

	// wrapper will add pid to procs file

	go func() {
		tracker.SetStatus(Running, 0)
		werr := cmd.Wait()
		if werr != nil {
			var ec int
			if ee, ok := werr.(*exec.ExitError); !ok {
				ec = ee.ExitCode()
			}
			tracker.SetStatus(Error, ec)
		} else {
			tracker.SetStatus(Completed, 0)
		}
		log.Println("... process ended")
		done <- true // shut down any streams
		ra.Cleanup()
	}()

	return cmd, tracker, done, nil
}
