package backup

import (
	rms_backup "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-backup"
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
	Size      uint64
	Type      rms_backup.BackupType
	Errors    []error
}

type state struct {
	mu sync.RWMutex
	r  Report
}

func (s *state) setInProgress(backupType rms_backup.BackupType, fileName string, createdAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.reset()

	s.r.Status = InProgress
	s.r.FileName = fileName
	s.r.Timestamp = createdAt
	s.r.Type = backupType
}

func (s *state) setFailed() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.r.Status = Failed
}

func (s *state) addError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.r.Errors = append(s.r.Errors, err)
}

func (s *state) setProgress(complete, total int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.r.Progress = (float32(complete) / float32(total)) * 100.
}

func (s *state) setReady(fileSize uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.r.Errors) != 0 {
		s.r.Status = ReadyWithErrors
	} else {
		s.r.Status = Ready
	}
	s.r.Size = fileSize
}

func (s *state) reset() {
	s.r.Size = 0
	s.r.Status = NeverRun
	s.r.FileName = ""
	s.r.Errors = nil
	s.r.Progress = 0
	s.r.Timestamp = time.Time{}
}

func (s *state) getReport() Report {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.r
}
