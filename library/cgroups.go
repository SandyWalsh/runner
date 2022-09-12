package library

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

func makeCGroups(cg *ControlGroup, ppid Process) {
	// makes cgroups in the v2 style of /sys/fs/cgroup/<ppid>
	group := fmt.Sprintf("/sys/fs/cgroup/%s", ppid)
	if err := os.MkdirAll(group, 0644); err != nil {
		if errors.Is(err, fs.ErrExist) {
			log.Println("cgroup", group, "already exists")
		} else {
			log.Println("failed to create cgroup for", group, ":", err)
			os.Exit(1)
		}
	}
	fn := fmt.Sprintf("%s/cpuset.cpus", group)
	alterController("1", fn)

	// Now we can write the actual cgroup limits ...
	for _, l := range cg.Limits {
		lf := fmt.Sprintf("%s/%s", group, l.Var)
		setVar(l.Value, lf)
	}
}

func alterController(cmd, fn string) {
	writeWithFlags(cmd, fn, os.O_WRONLY|os.O_APPEND)
}

func setVar(cmd, fn string) {
	writeWithFlags(cmd, fn, os.O_CREATE|os.O_WRONLY)
}

func writeWithFlags(cmd, fn string, flag int) {
	file, err := os.OpenFile(fn, flag, 0644)
	if err != nil {
		log.Println("unable to open", fn, "for writing :", err)
		os.Exit(1)
	}
	defer file.Close()

	if _, err := file.WriteString(cmd); err != nil {
		log.Println("error writing", cmd, "to", fn, ": ", err)
		os.Exit(1)
	}
	log.Println(cmd, ">>", fn)
}

func cgroupSetup(pid int) {
	//addProcessToCgroup(cpath+"cgroup.procs", pid)
}

type StatusTracker struct {
	Status   Status
	ExitCode int
	mtx      sync.Mutex
}

func (s *StatusTracker) GetStatus() (Status, int) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.Status, s.ExitCode
}

func (s *StatusTracker) SetStatus(st Status, ec int) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	// if we aborted ignore all subsequent attempts to set status
	if s.Status == Aborted {
		return
	}
	s.Status = st
	s.ExitCode = ec
}

func runCommand(ctx context.Context, wrapper string, spa *syscall.SysProcAttr, out *safeBuffer, cleanup func(), fn string, c string, args ...string) (*exec.Cmd, *StatusTracker, <-chan bool, error) {
	args = append([]string{fn, c}, args...)
	log.Println("running process:", wrapper, c, args)
	cmd := exec.Command(wrapper, args...)

	cmd.Stdout = out
	cmd.Stderr = out

	done := make(chan bool)

	cmd.SysProcAttr = spa // linux namespace controls

	tracker := &StatusTracker{Status: Unavailable}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error running the exec.Command - %s\n", err)
		cleanup()
		os.Exit(1)
	}

	// NOTE: wrapper will add pid to procs file

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
		done <- true
		cleanup()
	}()

	return cmd, tracker, done, nil
}
