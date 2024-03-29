package gitea

import (
	"fmt"
	"github.com/RacoonMediaServer/rms-backup/internal/backup"
	"github.com/RacoonMediaServer/rms-backup/internal/config"
	"github.com/RacoonMediaServer/rms-backup/internal/system"
	"os"
	"path/filepath"
)

const backupArtifact = "gitea.zip.bak"

type backupCommand struct {
	container string
}

func (b backupCommand) Title() string {
	return "Backup ALL"
}

func (b backupCommand) Execute(ctx backup.Context) error {
	const tmpPath = "/data/git/gitea-dump.zip"
	id, err := system.DockerGetContainerID(ctx, b.container)
	if err != nil {
		return fmt.Errorf("get container ID failed: %w", err)
	}
	if err = system.DockerExec(ctx, "git", id, "gitea", "dump", "-c", "/data/gitea/conf/app.ini", "-f", tmpPath); err != nil {
		return err
	}
	return system.DockerExec(ctx, "root", id, "mv", tmpPath, b.outputFile())
}

func (b backupCommand) outputFile() string {
	return filepath.Join(config.Config().Directories.Artifacts, backupArtifact)
}

func (b backupCommand) Cleanup(ctx backup.Context) error {
	return os.Remove(b.outputFile())
}
