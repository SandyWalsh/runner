package sprinter

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Runner allows for running linux processes safely.
type Runner interface {
	Run(ctx context.Context, client string, cmd string, args ...string) (Process, error)
	GetStatus(ctx context.Context, p Process) (Status, int, error)
	StreamOutput(ctx context.Context, p Process, sender Sender) error
	Abort(ctx context.Context, p Process) (Status, error)
}

// Process is a UUID which maps (internally) to a linux process ID
type Process string

// Sender allows the library to send stream info to the caller.
type Sender interface {
	Send(string) error
}

type runnerImpl struct {
	running     map[Process]*runningProcess
	authz       AuthZRules
	wrapperPath string
	mtx         sync.Mutex
	driver      Driver
}

func NewRunner(a AuthZRules, d Driver) (Runner, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	wrapper := filepath.Join(cwd, "../wrapper/wrapper")

	return &runnerImpl{
		running:     map[Process]*runningProcess{},
		wrapperPath: wrapper,
		driver:      d,
		authz:       a,
	}, nil
}

var _ Runner = (*runnerImpl)(nil)

func (r *runnerImpl) Run(ctx context.Context, client string, cmd string, args ...string) (Process, error) {
	cgn, ok := r.authz.ClientToCGroup[client]
	if !ok {
		return "", errors.New("authorization failed - no cgroup defined for this client")
	}
	cg, err := r.authz.getControlGroup(cgn)
	if err != nil {
		return "", err
	}

	puuid := Process(uuid.New().String())
	// make a temp dir for the process to run in. Keep it away from the server code.
	dir, err := os.MkdirTemp("", string(puuid)+"_*")
	if err != nil {
		return "", err
	}

	cgdir := filepath.Join("/sys/fs/cgroup/", string(puuid))
	if err := r.driver.MakeCgroup(cgdir, *cg); err != nil {
		return "", err
	}

	var buf SafeBuffer

	ra := RunArgs{
		Dir:              dir,
		WrapperPath:      r.wrapperPath,
		SPA:              cg.SysProcAttr,
		Buffer:           &buf,
		Cleanup:          cleanup(r, puuid), // the driver doesn't know about these things. Nor should it.
		ControlGroupPath: cgdir,
		Cmd:              cmd,
		Args:             args,
	}

	log.Println("wrapper", ra.WrapperPath, "cgpath", ra.ControlGroupPath, "cmd", ra.Cmd, "args", ra.Args)

	ps, st, done, err := r.driver.Run(ra)
	if err != nil {
		r.internalCleanup(puuid)
		return "", err
	}

	if ps.Process != nil {
		r.addRunningProcess(puuid, &runningProcess{cmd: ps, output: &buf, tracker: st, tempdir: dir, done: done})
		return puuid, nil
	}
	r.internalCleanup(puuid)
	return "", errors.New("process did not launch")
}

func (r *runnerImpl) GetStatus(ctx context.Context, p Process) (Status, int, error) {
	var rp *runningProcess
	var err error
	if rp, err = r.getRunningProcess(p); err != nil {
		return Unavailable, -1, err
	}
	st, ec := rp.tracker.GetStatus()
	return st, ec, nil
}

func (r *runnerImpl) StreamOutput(ctx context.Context, p Process, sender Sender) error {
	var rp *runningProcess
	var err error
	if rp, err = r.getRunningProcess(p); err != nil {
		return err
	}

	if st, _ := rp.tracker.GetStatus(); st != Running {
		return errors.New("process is not running")
	}

	rp.sender = sender // NOTE: CAN'T STORE IN RP

	log.Println("starting stdout stream")

	// NOTE: we can't do this in a goroutine.
	// If we return from this method the stream will be closed.
	done := false
	reader := rp.output.NewReader()
	for {
		done = true
		select {
		// look for signal when the process ends so we can bail out of here
		// and not require the client to do any status polling.
		//case <-rp.done:
		//	log.Println("process finished")
		//	done = true
		case <-ctx.Done():
			log.Println("caller context Done")
		default:
			done = false
			b := make([]byte, 256)
			n, rerr := reader.Read(b)
			log.Println("READ", n, rerr, string(b[:n]))
			if n > 0 {
				log.Println("SENDING:", string(b[:n]))
				if serr := sender.Send(string(b[:n])); serr != nil {
					log.Println("send error:", serr.Error())
					done = true
				}
			}

			if rerr != nil && rerr != io.EOF {
				done = true
				log.Println("read error", rerr)
			}
			time.Sleep(time.Second)
		}
		if done {
			break
		}
	}
	log.Println("closing output stream")
	return nil
}

func (r *runnerImpl) Abort(ctx context.Context, p Process) (Status, error) {
	var rp *runningProcess
	var err error
	if rp, err = r.getRunningProcess(p); err != nil {
		return Unavailable, err
	}
	rp.tracker.SetStatus(Aborted, 0)
	if err := rp.cmd.Process.Kill(); err != nil {
		log.Println("failed to kill process: ", err)
		return Unavailable, nil
	}
	return Aborted, nil
}

// ----

type runningProcess struct {
	cmd     *exec.Cmd
	sender  Sender
	output  *SafeBuffer
	tracker *StatusTracker
	done    <-chan bool
	tempdir string
}

func (r *runnerImpl) internalCleanup(puuid Process) {
	log.Println("looking up process", puuid)
	fn := filepath.Join("/sys/fs/cgroup/", string(puuid))
	if err := os.RemoveAll(fn); err != nil {
		log.Println("cannot clean up cgroup dir", fn, err)
	}

	// clean up temp dir
	var rp *runningProcess
	var err error
	if rp, err = r.getRunningProcess(puuid); err != nil {
		log.Println(err)
		return
	}

	log.Println("cleaning up", rp.tempdir)
	if err := os.RemoveAll(rp.tempdir); err != nil {
		log.Println("cannot clean up temp dir", rp.tempdir, err)
	}

	r.removeRunningProcess(puuid)
}

func cleanup(r *runnerImpl, puuid Process) func() {
	return func() {
		r.internalCleanup(puuid)
	}
}

func (r *runnerImpl) addRunningProcess(puuid Process, rp *runningProcess) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.running[puuid] = rp
}

func (r *runnerImpl) getRunningProcess(puuid Process) (*runningProcess, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	rp, ok := r.running[puuid]
	if !ok {
		return nil, fmt.Errorf("could not find process %s", puuid)
	}
	return rp, nil
}

func (r *runnerImpl) removeRunningProcess(puuid Process) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	delete(r.running, puuid)
}
