package system

import (
	"fmt"
	"github.com/RacoonMediaServer/rms-backup/internal/backup"
	"strings"
)

func DockerGetContainerID(ctx backup.Context, name string) (id string, err error) {
	var out string
	out, err = Exec(ctx, "docker", "ps", "-f", fmt.Sprintf("name=%s", name), "--format", "{{.ID}}")
	if err != nil {
		return
	}

	ids := strings.Split(out, "\n")
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
	args = append(args, command)
	args = append(args, params...)
	_, err := Exec(ctx, "docker", args...)
	return err
}
