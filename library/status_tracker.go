package library

import "sync"

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
