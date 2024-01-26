package nextcloud

import (
	"fmt"
	"github.com/RacoonMediaServer/rms-backup/internal/backup"
	"github.com/RacoonMediaServer/rms-backup/internal/system"
)

type setMaintenanceMode struct {
	name string
	id   string
}

func (c *setMaintenanceMode) Title() string {
	return "Enable/Disable maintenance mode"
}

func (c *setMaintenanceMode) Execute(ctx backup.Context) error {
	id, err := system.DockerGetContainerID(ctx, c.name)
	if err != nil {
		return fmt.Errorf("get container ID failed: %w", err)
	}
	c.id = id
	return system.DockerExec(ctx, "www-data", id, "php", "./occ", "maintenance:mode", "--on")
}

func (c *setMaintenanceMode) Cleanup(ctx backup.Context) error {
	return system.DockerExec(ctx, "www-data", c.id, "php", "./occ", "maintenance:mode", "--off")
}
