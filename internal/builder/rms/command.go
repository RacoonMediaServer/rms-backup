package rms

import (
	"fmt"
	"github.com/RacoonMediaServer/rms-backup/internal/backup"
	"github.com/RacoonMediaServer/rms-backup/internal/config"
	"github.com/RacoonMediaServer/rms-backup/internal/system"
	"github.com/RacoonMediaServer/rms-packages/pkg/configuration"
	"os"
	"path/filepath"
)

type backupDatabaseCmd struct {
	config configuration.Database
}

func (b backupDatabaseCmd) Title() string {
	return "Backup Settings"
}

func (b backupDatabaseCmd) Execute(ctx backup.Context) error {
	id, err := system.DockerGetContainerID(ctx, b.config.Host)
	if err != nil {
		return fmt.Errorf("get container ID failed: %w", err)
	}
	return system.DockerExec(ctx, "root", id, "pg_dump", b.config.Database, "-U", b.config.User, "-f", b.outputFile())
}

func (b backupDatabaseCmd) outputFile() string {
	return filepath.Join(config.Config().Directories.Artifacts, dbBackupArtifact)
}

func (b backupDatabaseCmd) Cleanup(ctx backup.Context) error {
	return os.Remove(b.outputFile())
}
