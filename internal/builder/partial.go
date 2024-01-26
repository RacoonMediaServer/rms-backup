package builder

import (
	"github.com/RacoonMediaServer/rms-backup/internal/backup"
	"github.com/RacoonMediaServer/rms-backup/internal/builder/nextcloud"
	"github.com/RacoonMediaServer/rms-backup/internal/config"
)

func createPartialBackup() backup.Instruction {
	i := backup.Instruction{Title: "Partial Backup"}
	i.Add(nextcloud.GetBackupStage(config.Config().Services, false))
	return i
}
