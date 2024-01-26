package builder

import (
	"github.com/RacoonMediaServer/rms-backup/internal/backup"
	"github.com/RacoonMediaServer/rms-backup/internal/builder/nextcloud"
	"github.com/RacoonMediaServer/rms-backup/internal/config"
)

func createFullBackup() backup.Instruction {
	i := backup.Instruction{Title: "Full Backup"}
	i.Add(nextcloud.GetBackupStage(config.Config().Services, true))
	return i
}
