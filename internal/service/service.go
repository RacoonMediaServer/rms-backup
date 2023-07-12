package service

import (
	"context"
	rms_backup "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-backup"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
}

func (s Service) LaunchBackup(ctx context.Context, request *rms_backup.LaunchBackupRequest, response *rms_backup.LaunchBackupResponse) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetBackupStatus(ctx context.Context, empty *emptypb.Empty, response *rms_backup.GetBackupStatusResponse) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetBackups(ctx context.Context, empty *emptypb.Empty, response *rms_backup.GetBackupsResponse) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) RemoveBackup(ctx context.Context, request *rms_backup.RemoveBackupRequest, empty *emptypb.Empty) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetBackupSettings(ctx context.Context, empty *emptypb.Empty, settings *rms_backup.BackupSettings) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) SetBackupSettings(ctx context.Context, settings *rms_backup.BackupSettings, empty *emptypb.Empty) error {
	//TODO implement me
	panic("implement me")
}

func NewService(interface{}) *Service {
	return &Service{}
}
