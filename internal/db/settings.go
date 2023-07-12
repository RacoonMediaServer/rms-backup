package db

import (
	"github.com/RacoonMediaServer/rms-backup/internal/config"
	rms_backup "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-backup"
)

type settings struct {
	ID       uint                      `gorm:"primaryKey"`
	Settings rms_backup.BackupSettings `gorm:"embedded"`
}

func (d *Database) LoadSettings() (*rms_backup.BackupSettings, error) {
	var record settings
	defaultSettings := settings{
		ID:       1,
		Settings: config.DefaultSettings,
	}
	if err := d.conn.Where(settings{ID: 1}).Attrs(defaultSettings).FirstOrCreate(&record).Error; err != nil {
		return nil, err
	}
	return &record.Settings, nil
}

func (d *Database) SaveSettings(val *rms_backup.BackupSettings) error {
	return d.conn.Save(&settings{ID: 1, Settings: *val}).Error
}
