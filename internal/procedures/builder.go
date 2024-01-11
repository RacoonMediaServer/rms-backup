package procedures

import (
	"github.com/RacoonMediaServer/rms-backup/internal/backup"
	"github.com/RacoonMediaServer/rms-backup/internal/procedures/nextcloud"
	rms_backup "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-backup"
)

func CreateBackupProcedure(backupType rms_backup.BackupType) backup.Instruction {
	switch backupType {
	case rms_backup.BackupType_Full:
		return createFullBackupProcedure()
	case rms_backup.BackupType_Partial:
		return createPartialBackupProcedure()
	default:
		panic("unknown backup type")
	}
}

func createFullBackupProcedure() backup.Instruction {
	p := backup.Instruction{Title: "FullBackup"}
	p.Add(nextcloud.GetBackupStage(true))
	return p
}

func createPartialBackupProcedure() backup.Instruction {
	p := backup.Instruction{Title: "PartialBackup"}
	p.Add(nextcloud.GetBackupStage(false))
	return p
}
