package sprinter_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"testing"

	"github.com/SandyWalsh/runner/sprinter"
)

type fakeDriver struct{}

func (f *fakeDriver) MakeCgroup(dir string, cg sprinter.ControlGroup) error {
	return nil
}

func (f *fakeDriver) Run(ra sprinter.RunArgs) (*exec.Cmd, *sprinter.StatusTracker, error) {
	st := &sprinter.StatusTracker{}
	e := &exec.Cmd{Process: &os.Process{Pid: 1}}
	return e, st, nil
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

	fcg := &fakeDriver{}
	lr, err := sprinter.NewRunner(authz, fcg)
	if err != nil {
		t.Fatal(err)
	}
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

func TestSafeBuffer(t *testing.T) {
	s := sprinter.SafeBuffer{}
	var msg string
	var wg sync.WaitGroup
	var writerErr error
	wg.Add(2)
	go func() {
		for i := 0; i < 10; i++ {
			s.Write([]byte("testing"))
		}
		s.Close()
		wg.Done()
	}()
	go func() {
		r, nd := s.NewReader()
		for {
			select {
			case x := <-nd:
				if x == 0 {
					goto done
				}
				b := make([]byte, x)
				n, err := r.Read(b)
				if n == 0 && err == io.EOF {
					goto done
				}
				if err != nil {
					writerErr = fmt.Errorf("got read error: %s", err)
					goto done
				}
				msg += string(b)
			}
		}
	done:
		wg.Done()
	}()
	wg.Wait()
	if writerErr != nil {
		t.Fatal(writerErr)
	}
	if len(msg) != len("testing")*10 {
		t.Fatal("expected 'testing' 10 times, got", msg)
	}
}
