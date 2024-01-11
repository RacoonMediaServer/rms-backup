package backup

import rms_backup "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-backup"

type InstructionBuilder interface {
	BuildBackup(backupType rms_backup.BackupType) Instruction
}
