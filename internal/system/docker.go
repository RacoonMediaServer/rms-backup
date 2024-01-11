package system

import (
	"context"
	"fmt"
	"github.com/RacoonMediaServer/rms-backup/internal/backup"
	"go-micro.dev/v4/logger"
	"os/exec"
	"strings"
)

func DockerGetContainerID(ctx context.Context, name string) (id string, err error) {
	var out []byte
	cmd := exec.CommandContext(ctx, "docker", "ps", "-f", fmt.Sprintf("name=%s", o.Name), "--format", "{{.ID}}")
	out, err = cmd.Output()
	if err != nil {
		return
	}

	ids := strings.Split(string(out), "\n")
	if len(ids) == 0 || len(ids[0]) == 0 {
		err = fmt.Errorf("container %s not found", name)
	}
	id = ids[0]

	return
}

func DockerExec(ctx backup.Context, user string, container string, command string, params ...string) error {
	args := []string{"exec"}
	if user != "" {
		args = append(args, "-u")
		args = append(args, user)
	}
	args = append(args, container)
	args = append(args, params...)
	cmd := exec.CommandContext(ctx, "docker", args...)
	out, err := cmd.Output()
	if err != nil {
		ctx.Log().Logf(logger.DebugLevel, "output:\n%s", string(out))
	}
	return err
}
