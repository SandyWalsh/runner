package sprinter_test

import (
	"context"
	"log"
	"os"
	"os/exec"
	"syscall"
	"testing"

	"github.com/SandyWalsh/runner/sprinter"
)

type fakeCgroup struct{}

func (f *fakeCgroup) MakeCgroup(cg *sprinter.ControlGroup, ppid sprinter.Process) error {
	return nil
}

func (f *fakeCgroup) Run(ctx context.Context, wrapper string, spa *syscall.SysProcAttr, out *sprinter.SafeBuffer, cleanup func(), cgDir string, c string, args ...string) (*exec.Cmd, *library.StatusTracker, <-chan bool, error) {
	done := make(chan bool)
	st := &sprinter.StatusTracker{Status: sprinter.Unavailable}
	e := &exec.Cmd{Process: &os.Process{Pid: 1}}
	return e, st, done, nil
}

func TestAuth(t *testing.T) {
	cg := sprinter.ControlGroup{
		Name: "test",
		Limits: []sprinter.Limit{
			{Var: "cpu.max", Value: "100000 1000000"},
			{Var: "cpu.weight", Value: "50"},
		}}
	authz := sprinter.AuthZRules{
		ControlGroups:  []sprinter.ControlGroup{cg},
		ClientToCGroup: map[string]string{"caller": "test"},
	}

	fcg := &fakeCgroup{}
	lr := sprinter.NewRunner(authz, fcg)
	ctx := context.Background()
	pid, err := lr.Run(ctx, "bad_actor", "ls", "-la")
	if err == nil {
		t.Fatal("expected auth err, got none")
	}
	pid, err = lr.Run(ctx, "caller", "ls", "-la")
	if err != nil {
		t.Fatal("expected pid, got err: ", err)
	}
	st, ec, err := lr.GetStatus(ctx, pid)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("status", st, ec)
}
