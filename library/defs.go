package library

import (
	"context"
	"syscall"
)

// Process is a UUID which maps (internally) to a linux process ID
type Process string

// Status reflects the status of a running process
type Status int64

const (
	Unavailable Status = 0
	Error              = 1
	Running            = 2
	Completed          = 3
	Aborted            = 4
)

// Sender allows the library to send stream info to the caller.
type Sender interface {
	Send(string) error
	GetContext() context.Context // Is there a context we should be watching for Done?
}

// Limit is a single cgroup v2 parameter and value
type Limit struct {
	Var   string
	Value string
}

// ControlGroup defines a cgroups v2 namespace and related resource limits
type ControlGroup struct {
	Name        string
	Limits      []Limit
	SysProcAttr *syscall.SysProcAttr
}

// AuthZRules defines a set of cgroups and the mapping of a client id to these cgroups when running processes.
type AuthZRules struct {
	ControlGroups  []ControlGroup
	ClientToCGroup map[string]string
}

// Runner allows for running linux processes safely.
type Runner interface {
	Run(ctx context.Context, client string, cmd string, args ...string) (Process, error)
	GetStatus(ctx context.Context, p Process) (Status, int, error)
	StreamOutput(ctx context.Context, p Process, sender Sender) error
	Abort(ctx context.Context, p Process) (Status, error)
}
