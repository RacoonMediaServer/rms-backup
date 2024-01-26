package service

import (
	"github.com/RacoonMediaServer/rms-backup/internal/config"
	"path/filepath"
)

func getAbsoluteFileName(fileName string) string {
	return filepath.Join(config.Config().Directories.Backups, fileName)
}
