package backup

import (
	"fmt"
	rms_backup "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-backup"
	"time"
)

const backupOutputExtension = "7z"

func genFileName(backupType rms_backup.BackupType, createdAt time.Time) string {
	return fmt.Sprintf("%s_%s.%s", backupType.String(), createdAt.Format(time.DateOnly), backupOutputExtension)
}
