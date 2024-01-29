package service

import (
	"context"
	"fmt"
	"github.com/RacoonMediaServer/rms-backup/internal/backup"
	"github.com/RacoonMediaServer/rms-backup/internal/model"
	"github.com/RacoonMediaServer/rms-packages/pkg/events"
	"github.com/RacoonMediaServer/rms-packages/pkg/misc"
	rms_backup "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-backup"
	"go-micro.dev/v4/logger"
	"time"
)

const publishTimeout = 40 * time.Second

func (s *Service) startRegularBackup() {
	logger.Info("Start regular backup...")

	s.mu.RLock()
	backupType := s.settings.Type
	s.mu.RUnlock()

	s.startBackup(backupType)
}

func (s *Service) startBackup(backupType rms_backup.BackupType) bool {
	instructions := s.builder.BuildBackup(backupType)
	return s.engine.Launch(backup.NewContext(), backupType, instructions)
}

func (s *Service) onBackupReady(report backup.Report) {
	switch report.Status {
	case backup.Ready:
		s.storeBackupInfo(report)
		s.notifyBackupDone(report)
	case backup.Failed:
		s.notifyBackupFailed(report)
	case backup.ReadyWithErrors:
		s.storeBackupInfo(report)
		s.notifyBackupFailed(report)
		s.notifyBackupDone(report)
	default:
		logger.Errorf("Unknown backup status: %d", report.Status)
	}
}

func (s *Service) notifyBackupDone(report backup.Report) {
	sizeMB := uint32(report.Size / (1024 * 1024))
	event := events.Notification{
		Sender:    "rms-backup",
		Kind:      events.Notification_BackupComplete,
		ItemTitle: &report.FileName,
		SizeMB:    &sizeMB,
	}
	s.publishEvent(&event)
	logger.Infof("Backup data saved to '%s'", report.FileName)
}

func (s *Service) notifyBackupFailed(report backup.Report) {
	logger.Errorf("Backup process failed: %+v", report.Errors)
	event := events.Malfunction{
		Sender:     "rms-backup",
		Timestamp:  time.Now().Unix(),
		Error:      fmt.Sprintf("Возникли %d ошибок при резервном копировании:\n%+v", len(report.Errors), report.Errors),
		System:     events.Malfunction_Services,
		Code:       events.Malfunction_ActionFailed,
		StackTrace: misc.GetStackTrace(),
	}
	s.publishEvent(&event)
}

func (s *Service) publishEvent(event interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), publishTimeout)
	defer cancel()

	if err := s.pub.Publish(ctx, event); err != nil {
		logger.Errorf("Publish event failed: %s", err)
	}
}

func (s *Service) storeBackupInfo(report backup.Report) {
	br := model.BackupRecord{
		FileName:  report.FileName,
		CreatedAt: report.Timestamp,
		Type:      report.Type,
		Size:      report.Size,
	}
	if err := s.db.AddBackup(&br); err != nil {
		logger.Errorf("Add backup info to database failed: %s", err)
	}
}
