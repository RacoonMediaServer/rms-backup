package nextcloud

import "github.com/RacoonMediaServer/rms-backup/internal/backup"

func GetBackupStage(includeData bool) backup.Stage {
	s := backup.Stage{Title: "Backup Nextcloud"}
	s.Add(&setMaintenanceMode{})
	return s
}
