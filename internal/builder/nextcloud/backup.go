package nextcloud

import (
	"github.com/RacoonMediaServer/rms-backup/internal/backup"
	"github.com/RacoonMediaServer/rms-backup/internal/config"
)

func GetBackupStage(services config.Services, includeData bool) backup.Stage {
	s := backup.Stage{Title: "Backup Nextcloud"}
	s.Add(&setMaintenanceMode{name: services.Nextcloud.Container})
	s.Add(&dbBackupCommand{dbContainer: services.Postgres.Container, nextcloud: services.Nextcloud})
	s.Artifacts = append(s.Artifacts, databaseBackupArtifact)
	if includeData {
		s.Artifacts = append(s.Artifacts, services.Nextcloud.Data)
	}
	return s
}
