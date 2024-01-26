package builder

import (
	"github.com/RacoonMediaServer/rms-backup/internal/backup"
	"github.com/RacoonMediaServer/rms-backup/internal/builder/gitea"
	"github.com/RacoonMediaServer/rms-backup/internal/builder/nextcloud"
	"github.com/RacoonMediaServer/rms-backup/internal/builder/rms"
	"github.com/RacoonMediaServer/rms-backup/internal/config"
)

func createPartialBackup() backup.Instruction {
	i := backup.Instruction{Title: "Partial Backup"}
	conf := config.Config()
	i.Add(nextcloud.GetBackupStage(conf.Services, false))
	i.Add(gitea.GetBackupStage(conf.Services))
	i.Add(rms.GetBackupStage(conf.Database))
	return i
}
