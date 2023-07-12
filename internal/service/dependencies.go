package service

import rms_backup "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-backup"

type Database interface {
	LoadSettings() (*rms_backup.BackupSettings, error)
}
