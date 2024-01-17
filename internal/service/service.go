package service

import (
	"context"
	"fmt"
	"github.com/RacoonMediaServer/rms-backup/internal/backup"
	"github.com/RacoonMediaServer/rms-backup/internal/compressor"
	"github.com/RacoonMediaServer/rms-backup/internal/config"
	"github.com/RacoonMediaServer/rms-packages/pkg/misc"
	rms_backup "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-backup"
	"github.com/go-co-op/gocron"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
	"sync"
	"time"
)

type Service struct {
	db Database

	sched   *gocron.Scheduler
	engine  *backup.Engine
	builder backup.InstructionBuilder
	pub     micro.Event

	mu       sync.RWMutex
	settings *rms_backup.BackupSettings
	job      *gocron.Job
}

func (s *Service) LaunchBackup(ctx context.Context, request *rms_backup.LaunchBackupRequest, response *rms_backup.LaunchBackupResponse) error {
	response.AlreadyLaunch = !s.startBackup(request.Type)
	return nil
}

func (s *Service) GetBackupStatus(ctx context.Context, empty *emptypb.Empty, response *rms_backup.GetBackupStatusResponse) error {
	r := s.engine.GetReport()
	response.Progress = r.Progress
	response.LastTime = uint64(r.Timestamp.Unix())
	switch r.Status {
	case backup.NeverRun:
		response.Status = rms_backup.GetBackupStatusResponse_NotStarted
	case backup.ReadyWithErrors:
		fallthrough
	case backup.Ready:
		fallthrough
	case backup.Failed:
		response.Status = rms_backup.GetBackupStatusResponse_Ready
	case backup.InProgress:
		response.Status = rms_backup.GetBackupStatusResponse_InProgress
	}

	return nil
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

func NewService(db Database, builder backup.InstructionBuilder, pub micro.Event) *Service {
	service := &Service{
		db:      db,
		sched:   gocron.NewScheduler(time.Local),
		builder: builder,
		pub:     pub,
	}

	compressionSettingsProvider := func() compressor.Settings {
		service.mu.RLock()
		defer service.mu.RUnlock()
		return compressor.Settings{Password: service.settings.Password}
	}
	compr := compressor.New(compressor.Format_7z, compressionSettingsProvider)

	engine := backup.NewEngine(compr)
	engine.SetTimeout(config.Config().BackupTimeout())
	engine.OnReady = service.onBackupReady

	service.engine = engine
	return service
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
