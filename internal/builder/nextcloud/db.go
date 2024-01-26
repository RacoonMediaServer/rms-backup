package nextcloud

import (
	"fmt"
	"github.com/RacoonMediaServer/rms-backup/internal/backup"
	"github.com/RacoonMediaServer/rms-backup/internal/config"
	"github.com/RacoonMediaServer/rms-backup/internal/system"
	"os"
	"path/filepath"
)

const databaseBackupArtifact = "nextcloud.db.bak"

type dbBackupCommand struct {
	dbContainer string
	nextcloud   config.Nextcloud
}

func (d dbBackupCommand) Title() string {
	return "Backup Database"
}

func (d dbBackupCommand) outputFile() string {
	return filepath.Join(d.nextcloud.InternalDirectory, databaseBackupArtifact)
}

func (d dbBackupCommand) Execute(ctx backup.Context) error {
	id, err := system.DockerGetContainerID(ctx, d.dbContainer)
	if err != nil {
		return fmt.Errorf("get container ID failed: %w", err)
	}
	return system.DockerExec(ctx, "root", id, "pg_dump", d.nextcloud.Database, "-U", d.nextcloud.User, "-f", d.outputFile())
}

func (d dbBackupCommand) Cleanup(ctx backup.Context) error {
	return os.Remove(d.outputFile())
}
