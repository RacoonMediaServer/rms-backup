package system

import (
	"github.com/RacoonMediaServer/rms-backup/internal/backup"
	"go-micro.dev/v4/logger"
	"os/exec"
)

func Exec(ctx backup.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	ctx.Log().Log(logger.DebugLevel, cmd.String())
	output, err := cmd.Output()
	if err != nil {
		ctx.Log().Logf(logger.DebugLevel, "Command output:\n%s", output)
	}
	return string(output), err
}
