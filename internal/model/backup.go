package model

import (
	rms_backup "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-backup"
	"time"
)

type BackupRecord struct {
	FileName  string `gorm:"primaryKey"`
	CreatedAt time.Time
	Type      rms_backup.BackupType
	Size      uint64
}

func (r BackupRecord) Convert() *rms_backup.BackupInfo {
	return &rms_backup.BackupInfo{
		FileName: r.FileName,
		Date:     uint64(r.CreatedAt.Unix()),
		Type:     r.Type,
		Size:     r.Size,
	}
}
