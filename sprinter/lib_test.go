package sprinter_test

import (
	"context"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/SandyWalsh/runner/sprinter"
)

type fakeCgroup struct{}

func (f *fakeCgroup) MakeCgroup(dir string, cg sprinter.ControlGroup) error {
	return nil
}

func (f *fakeCgroup) Run(ra sprinter.RunArgs) (*exec.Cmd, *sprinter.StatusTracker, <-chan bool, error) {
	done := make(chan bool)
	st := &sprinter.StatusTracker{}
	e := &exec.Cmd{Process: &os.Process{Pid: 1}}
	return e, st, done, nil
}

func TestAuth(t *testing.T) {
	cg := map[string]sprinter.ControlGroup{
		"test": {
			Limits: []sprinter.Limit{
				{Var: "cpu.max", Value: "100000 1000000"},
				{Var: "cpu.weight", Value: "50"},
			}}}
	authz := sprinter.AuthZRules{
		ControlGroups:  cg,
		ClientToCGroup: map[string]string{"caller": "test"},
	}

	fcg := &fakeCgroup{}
	lr, err := sprinter.NewRunner(authz, fcg)
	if err != nil {
		t.Fatal("unable to initialize sprinter library", err)
	}
	ctx := context.Background()
	pid, err := lr.Run(ctx, "bad_actor", "ls", "-la")
	if err == nil {
		t.Fatal("expected auth err, got none")
	}
	_, err = lr.Run(ctx, "caller", "ls", "-la")
	if err != nil {
		t.Fatal("expected pid, got err: ", err)
	}
	st, ec, err := lr.GetStatus(ctx, pid)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("status", st, ec)
}
