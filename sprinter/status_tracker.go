package sprinter

import "sync"

type StatusTracker struct {
	status   Status
	exitCode int
	mtx      sync.Mutex
}

// Status reflects the status of a running process
type Status int64

const (
	Unavailable Status = iota
	Error
	Running
	Completed
	Aborted
)

func (s *StatusTracker) GetStatus() (Status, int) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.status, s.exitCode
}

func (s *StatusTracker) SetStatus(st Status, ec int) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	// if we aborted ignore all subsequent attempts to set status
	if s.status == Aborted {
		return
	}
	s.status = st
	s.exitCode = ec
}
