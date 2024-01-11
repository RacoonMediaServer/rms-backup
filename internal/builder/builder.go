package builder

import (
	"github.com/RacoonMediaServer/rms-backup/internal/backup"
	rms_backup "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-backup"
)

type builder struct {
}

func (b builder) BuildBackup(backupType rms_backup.BackupType) backup.Instruction {
	switch backupType {
	case rms_backup.BackupType_Full:
		return createFullBackup()
	case rms_backup.BackupType_Partial:
		return createPartialBackup()
	default:
		panic("unknown backup type")
	}
}

func New() backup.InstructionBuilder {
	return &builder{}
}
