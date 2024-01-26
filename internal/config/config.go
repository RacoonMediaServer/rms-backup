package config

import (
	"github.com/RacoonMediaServer/rms-packages/pkg/configuration"
	"time"
)

// Configuration represents entire service configuration
type Configuration struct {
	BackupTimeoutSec int64 `json:"backup-timeout"`
	Database         configuration.Database
	Services         Services
	Directories      Directories
}

var config Configuration

// Load open and parses configuration file
func Load(configFilePath string) error {
	return configuration.Load(configFilePath, &config)
}

// Config returns loaded configuration
func Config() Configuration {
	return config
}

func (c Configuration) BackupTimeout() time.Duration {
	return time.Duration(c.BackupTimeoutSec) * time.Second
}
