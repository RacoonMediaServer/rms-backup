package rms

import (
	"github.com/RacoonMediaServer/rms-backup/internal/backup"
	"github.com/RacoonMediaServer/rms-packages/pkg/configuration"
)

const dbBackupArtifact = "rms.sql.bak"

func GetBackupStage(conf configuration.Database) backup.Stage {
	s := backup.Stage{Title: "Backup RMS settings"}
	s.Add(&backupDatabaseCmd{config: conf})
	s.Artifacts = []string{dbBackupArtifact}
	return s
}
