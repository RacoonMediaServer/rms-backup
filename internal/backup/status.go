package backup

import (
	"sync"
	"time"
)

type Status int

const (
	NeverRun Status = iota
	InProgress
	Failed
	ReadyWithErrors
	Ready
)

type Report struct {
	Status    Status
	FileName  string
	Timestamp time.Time
	Progress  float32
}

type state struct {
	mu sync.RWMutex
	r  Report
}

func (s *state) setInProgress() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.r.Status = InProgress
	s.r.FileName = ""
	s.r.Timestamp = time.Now()
}

func (s *state) setFailed() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.r.Status = Failed
}

func (s *state) setProgress(complete, total int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.r.Progress = (float32(complete) / float32(total)) * 100.
}

func (s *state) setReady(haveErrors bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if haveErrors {
		s.r.Status = ReadyWithErrors
	} else {
		s.r.Status = Ready
	}
}

func (s *state) getReport() Report {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.r
}
