package gitea

import (
	"github.com/RacoonMediaServer/rms-backup/internal/backup"
	"github.com/RacoonMediaServer/rms-backup/internal/config"
)

func GetBackupStage(services config.Services) backup.Stage {
	s := backup.Stage{Title: "Backup Gitea"}
	s.Add(&backupCommand{container: services.Gitea.Container})
	s.Artifacts = []string{backupArtifact}
	return s
}
