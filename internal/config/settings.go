package config

import rms_backup "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-backup"

var DefaultSettings = rms_backup.BackupSettings{
	Enabled:  true,
	Type:     rms_backup.BackupType_Full,
	Period:   rms_backup.BackupSettings_EveryMonth,
	Day:      1,
	Hour:     0,
	Password: nil,
}
