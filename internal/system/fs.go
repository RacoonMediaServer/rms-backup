package system

import (
	"github.com/RacoonMediaServer/rms-backup/internal/backup"
	"go-micro.dev/v4/logger"
	"os/exec"
)

func RsyncCopy(ctx backup.Context, src, dst string) error {
	cmd := exec.CommandContext(ctx, "rsync", "-Aax", src, dst)
	out, err := cmd.Output()
	if err != nil {
		ctx.Log().Logf(logger.DebugLevel, "rsync output:\n%s", string(out))
		return err
	}
	return nil
}
