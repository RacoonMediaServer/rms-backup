package builder

import (
	"github.com/RacoonMediaServer/rms-backup/internal/backup"
	"github.com/RacoonMediaServer/rms-backup/internal/builder/gitea"
	"github.com/RacoonMediaServer/rms-backup/internal/builder/nextcloud"
	"github.com/RacoonMediaServer/rms-backup/internal/builder/rms"
	"github.com/RacoonMediaServer/rms-backup/internal/config"
)

func createFullBackup() backup.Instruction {
	i := backup.Instruction{Title: "Full Backup"}
	conf := config.Config()
	i.Add(nextcloud.GetBackupStage(conf.Services, true))
	i.Add(gitea.GetBackupStage(conf.Services))
	i.Add(rms.GetBackupStage(conf.Database))
	return i
}
