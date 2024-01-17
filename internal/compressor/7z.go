package compressor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

type zip7 struct {
	provider SettingsProvider
}

func (z zip7) Extension() string {
	return "7z"
}

func (z zip7) Compress(ctx context.Context, targetArchive string, files []string) (uint64, error) {
	args := []string{"a", targetArchive}
	if z.provider != nil {
		settings := z.provider()
		if settings.Password != nil {
			args = append(args, fmt.Sprintf("-p%s", *settings.Password), "-mhe")
		}
	}
	args = append(args, files...)
	cmd := exec.CommandContext(ctx, "7z", args...)
	output, err := cmd.Output()
	if err != nil {
		err = fmt.Errorf("%w:\n%s", err, string(output))
		return 0, err
	}

	var fi os.FileInfo
	fi, err = os.Stat(targetArchive)
	if err != nil {
		return 0, err
	}
	return uint64(fi.Size()), err
}
