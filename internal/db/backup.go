package db

import "github.com/RacoonMediaServer/rms-backup/internal/model"

func (d *Database) LoadBackups() (result []model.BackupRecord, err error) {
	err = d.conn.Order("created_at asc").Find(&result).Error
	return
}

func (d *Database) AddBackup(backup *model.BackupRecord) error {
	return d.conn.Create(backup).Error
}

func (d *Database) RemoveBackup(fileName string) error {
	return d.conn.Model(&model.BackupRecord{}).Unscoped().Delete(&model.BackupRecord{FileName: fileName}).Error
}
