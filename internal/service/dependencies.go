package service

import (
	"github.com/RacoonMediaServer/rms-backup/internal/model"
	rms_backup "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-backup"
)

type Database interface {
	LoadSettings() (*rms_backup.BackupSettings, error)
	LoadBackups() (result []model.BackupRecord, err error)
	AddBackup(backup *model.BackupRecord) error
	RemoveBackup(fileName string) error
}

type Engine interface {
}
