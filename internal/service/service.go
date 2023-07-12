package service

import (
	"context"
	"fmt"
	rms_backup "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-backup"
	"github.com/go-co-op/gocron"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
	"time"
)

type Service struct {
	db Database

	sched *gocron.Scheduler

	mu       sync.RWMutex
	settings *rms_backup.BackupSettings
	job      *gocron.Job
}

func (s *Service) LaunchBackup(ctx context.Context, request *rms_backup.LaunchBackupRequest, response *rms_backup.LaunchBackupResponse) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetBackupStatus(ctx context.Context, empty *emptypb.Empty, response *rms_backup.GetBackupStatusResponse) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetBackups(ctx context.Context, empty *emptypb.Empty, response *rms_backup.GetBackupsResponse) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) RemoveBackup(ctx context.Context, request *rms_backup.RemoveBackupRequest, empty *emptypb.Empty) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetBackupSettings(ctx context.Context, empty *emptypb.Empty, settings *rms_backup.BackupSettings) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) SetBackupSettings(ctx context.Context, settings *rms_backup.BackupSettings, empty *emptypb.Empty) error {
	//TODO implement me
	panic("implement me")
}

func NewService(db Database) *Service {
	return &Service{
		db:    db,
		sched: gocron.NewScheduler(time.Local),
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
