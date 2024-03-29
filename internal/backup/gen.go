package backup

import (
	"fmt"
	rms_backup "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-backup"
	"time"
)

func (e *Engine) genFileName(backupType rms_backup.BackupType, createdAt time.Time) string {
	return fmt.Sprintf("Backup_%s_%s.%s", backupType.String(), createdAt.Format("2006-01-02T150405"), e.compressor.Extension())
}
