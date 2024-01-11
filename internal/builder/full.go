package builder

import (
	"github.com/RacoonMediaServer/rms-backup/internal/backup"
	"github.com/RacoonMediaServer/rms-backup/internal/builder/nextcloud"
)

func createFullBackup() backup.Instruction {
	p := backup.Instruction{Title: "FullBackup"}
	p.Add(nextcloud.GetBackupStage(true))
	return p
}
