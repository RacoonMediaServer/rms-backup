package service

import (
	"context"
	"fmt"
	"github.com/RacoonMediaServer/rms-backup/internal/backup"
	"github.com/RacoonMediaServer/rms-backup/internal/procedures"
	"github.com/RacoonMediaServer/rms-packages/pkg/misc"
	rms_backup "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-backup"
	"github.com/go-co-op/gocron"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
	"sync"
	"time"
)

type Service struct {
	db Database

	sched *gocron.Scheduler

	mu       sync.RWMutex
	settings *rms_backup.BackupSettings
	job      *gocron.Job

	engine *backup.Engine
}

func (s *Service) LaunchBackup(ctx context.Context, request *rms_backup.LaunchBackupRequest, response *rms_backup.LaunchBackupResponse) error {
	proc := procedures.CreateBackupProcedure(request.Type)
	response.AlreadyLaunch = !s.engine.Launch(backup.NewContext(), proc)
	return nil
}

func (s *Service) GetBackupStatus(ctx context.Context, empty *emptypb.Empty, response *rms_backup.GetBackupStatusResponse) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetBackups(ctx context.Context, empty *emptypb.Empty, response *rms_backup.GetBackupsResponse) error {
	list, err := s.db.LoadBackups()
	if err != nil {
		logger.Errorf("Load backups from database failed: %s", err)
		return err
	}
	response.Backups = make([]*rms_backup.BackupInfo, len(list))
	for i := range list {
		response.Backups[i] = list[i].Convert()
	}
	return nil
}

func (s *Service) RemoveBackup(ctx context.Context, request *rms_backup.RemoveBackupRequest, empty *emptypb.Empty) error {
	if err := s.db.RemoveBackup(request.FileName); err != nil {
		logger.Errorf("Remove backup '%s' from database failed: %s", request.FileName, err)
		return err
	}
	if err := os.Remove(getAbsoluteFileName(request.FileName)); err != nil {
		logger.Warnf("Remove backup '%s' failed: %s", request.FileName, err)
	}
	return nil
}

func (s *Service) GetBackupSettings(ctx context.Context, empty *emptypb.Empty, settings *rms_backup.BackupSettings) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	misc.AssignFields(s.settings, settings)
	return nil
}

func (s *Service) SetBackupSettings(ctx context.Context, settings *rms_backup.BackupSettings, empty *emptypb.Empty) error {
	logger.Infof("Change settings: %+v", settings)
	s.setSettings(settings)
	return nil
}

func NewService(db Database) *Service {
	return &Service{
		db:     db,
		sched:  gocron.NewScheduler(time.Local),
		engine: backup.NewEngine(),
	}
}

func (s *Service) Start() error {
	settings, err := s.db.LoadSettings()
	if err != nil {
		return fmt.Errorf("load settings failed: %w", err)
	}
	s.sched.StartAsync()

	s.setSettings(settings)
	return nil
}
