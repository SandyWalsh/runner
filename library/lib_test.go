package library_test

import (
	"context"
	"log"
	"os"
	"os/exec"
	"syscall"
	"testing"

	"github.com/SandyWalsh/runner/library"
)

type fakeCgroup struct{}

func (f *fakeCgroup) MakeCgroup(cg *library.ControlGroup, ppid library.Process) error {
	return nil
}

func (f *fakeCgroup) Run(ctx context.Context, wrapper string, spa *syscall.SysProcAttr, out *library.SafeBuffer, cleanup func(), cgDir string, c string, args ...string) (*exec.Cmd, *library.StatusTracker, <-chan bool, error) {
	done := make(chan bool)
	st := &library.StatusTracker{Status: library.Unavailable}
	e := &exec.Cmd{Process: &os.Process{Pid: 1}}
	return e, st, done, nil
}

func TestAuth(t *testing.T) {
	cg := library.ControlGroup{
		Name: "test",
		Limits: []library.Limit{
			{Var: "cpu.max", Value: "100000 1000000"},
			{Var: "cpu.weight", Value: "50"},
		}}
	authz := library.AuthZRules{
		ControlGroups:  []library.ControlGroup{cg},
		ClientToCGroup: map[string]string{"caller": "test"},
	}

	fcg := &fakeCgroup{}
	lr := library.NewRunner(authz, fcg)
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
