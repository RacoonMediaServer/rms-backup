package service

import (
	"fmt"
	rms_backup "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-backup"
	"github.com/go-co-op/gocron"
	"time"
)

func (s *Service) setSettings(settings *rms_backup.BackupSettings) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.settings = settings

	if s.job != nil {
		s.sched.RemoveByReference(s.job)
		s.job = nil
	}

	if !settings.Enabled {
		return
	}

	var sched *gocron.Scheduler
	switch settings.Period {
	case rms_backup.BackupSettings_EveryDay:
		sched = s.sched.Every(1).Day()
	case rms_backup.BackupSettings_EveryWeek:
		sched = s.sched.Every(1).Weekday(time.Weekday(settings.Day))
	case rms_backup.BackupSettings_EveryMonth:
		sched = s.sched.Every(1).Month(int(settings.Day))
	}
	job, err := sched.At(fmt.Sprintf("%d:00", settings.Hour)).Do(s.startRegularBackup)
	if err != nil {
		panic(err)
	}

	s.job = job
}
