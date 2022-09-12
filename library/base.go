package library

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"sync"

	"github.com/google/uuid"
)

type runningProcess struct {
	cmd     *exec.Cmd
	sender  Sender
	output  *safeBuffer
	tracker *StatusTracker
	done    <-chan bool
	pid     int
	tempdir string
}

type RunnerImpl struct {
	running     map[Process]*runningProcess
	authz       AuthZRules
	wrapperPath string
	mtx         sync.Mutex
}

func NewRunner(ns string, a AuthZRules) *RunnerImpl {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	wrapper := fmt.Sprintf("%s/../wrapper/wrapper", cwd)
	wrapper = path.Clean(wrapper)

	return &RunnerImpl{
		running:     map[Process]*runningProcess{},
		wrapperPath: wrapper,
		authz:       a,
	}
}

var _ Runner = (*RunnerImpl)(nil)

func (a *AuthZRules) getControlGroup(name string) (*ControlGroup, error) {
	for _, cg := range a.ControlGroups {
		if cg.Name == name {
			return &cg, nil
		}
	}
	return nil, errors.New("no such control group name")
}

func (r *RunnerImpl) internalCleanup(puuid Process) {
	log.Println("looking up process", puuid)
	fn := fmt.Sprintf("/sys/fs/cgroup/%s", puuid)
	if err := os.RemoveAll(fn); err != nil {
		log.Fatalln("cannot clean up cgroup dir", fn, err)
	}

	// clean up temp dir
	if rp, ok := r.getRunningProcess(puuid); ok {
		log.Println("cleaning up", rp.tempdir)
		if err := os.RemoveAll(rp.tempdir); err != nil {
			log.Fatalln("cannot clean up temp dir", rp.tempdir, err)
		}
	} else {
		log.Println("could not find process", puuid)
	}
}

func cleanup(r *RunnerImpl, puuid Process) func() {
	return func() {
		r.internalCleanup(puuid)
	}
}

func (r *RunnerImpl) addRunningProcess(puuid Process, rp *runningProcess) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.running[puuid] = rp
}

func (r *RunnerImpl) getRunningProcess(puuid Process) (*runningProcess, bool) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	rp, ok := r.running[puuid]
	return rp, ok
}

func (r *RunnerImpl) removeRunningProcess(puuid Process) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	delete(r.running, puuid)
}

func (r *RunnerImpl) Run(ctx context.Context, client string, cmd string, args ...string) (Process, error) {
	cgn, ok := r.authz.ClientToCGroup[client]
	if !ok {
		return "", errors.New("authorization failed - no cgroup defined for this client")
	}
	cg, err := r.authz.getControlGroup(cgn)
	if err != nil {
		return "", err
	}

	puuid := Process(uuid.New().String())

	dir, err := ioutil.TempDir("", string(puuid)+"_*")
	if err != nil {
		return "", err
	}

	os.Chdir(dir)

	makeCGroups(cg, puuid)

	fn := fmt.Sprintf("/sys/fs/cgroup/%s", puuid)

	var buf safeBuffer
	ps, st, done, err := runCommand(ctx, r.wrapperPath, cg.SysProcAttr, &buf, cleanup(r, puuid), fn, cmd, args...)
	if err != nil {
		r.internalCleanup(puuid)
		return "", err
	}

	if ps.Process != nil {
		r.addRunningProcess(puuid, &runningProcess{cmd: ps, output: &buf, tracker: st, pid: ps.Process.Pid, tempdir: dir, done: done})
		return puuid, nil
	}
	r.internalCleanup(puuid)
	return "", errors.New("process did not launch")
}

func (r *RunnerImpl) safeGetRunningProcess(p Process) (*runningProcess, error) {
	rp, ok := r.running[p]
	if !ok {
		return nil, errors.New("process not found")
	}
	return rp, nil
}

func (r *RunnerImpl) GetStatus(ctx context.Context, p Process) (Status, int, error) {
	rp, err := r.safeGetRunningProcess(p)
	if err != nil {
		return Unavailable, -1, err
	}
	st, ec := rp.tracker.GetStatus()
	return st, ec, nil
}

func (r *RunnerImpl) StreamOutput(ctx context.Context, p Process, sender Sender) error {
	rp, err := r.safeGetRunningProcess(p)
	if err != nil {
		return err
	}

	if st, _ := rp.tracker.GetStatus(); st != Running {
		return errors.New("process is not running")
	}

	rp.sender = sender

	log.Println("starting stdout stream")

	// NOTE: we can't do this in a goroutine.
	// If we return from this method the stream will be closed.
	done := false
	reader := rp.output.NewReader()
	for {
		select {
		// look for signal when the process ends so we can bail out of here
		// and not require the client to do any status polling.
		case <-rp.done:
			log.Println("process finished")
			done = true
		case <-sender.GetContext().Done():
			log.Println("context Done")
			done = true
		default:
			b := make([]byte, 256)
			_, err := reader.Read(b)
			if err != nil {
				// We'll get EOF if we're caught up, but that doesn't mean the process has finished.
				if err != io.EOF {
					log.Println("reader error", err)
					break
				}
			}
			err = sender.Send(string(b))
			if err != nil {
				log.Println("send error:", err.Error())
			}
		}
		if done {
			break
		}
	}
	log.Println("closing output stream")

	return nil
}

func (r *RunnerImpl) Abort(ctx context.Context, p Process) (Status, error) {
	rp, err := r.safeGetRunningProcess(p)
	if err != nil {
		return Unavailable, err
	}
	rp.tracker.SetStatus(Aborted, 0)
	if err := rp.cmd.Process.Kill(); err != nil {
		log.Fatal("failed to kill process: ", err)
		return Unavailable, nil
	}
	return Aborted, nil
}
