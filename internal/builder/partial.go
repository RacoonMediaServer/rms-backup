package builder

import (
	"github.com/RacoonMediaServer/rms-backup/internal/backup"
	"github.com/RacoonMediaServer/rms-backup/internal/builder/nextcloud"
)

func createPartialBackup() backup.Instruction {
	p := backup.Instruction{Title: "PartialBackup"}
	p.Add(nextcloud.GetBackupStage(false))
	return p
}
