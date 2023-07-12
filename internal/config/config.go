package config

import "github.com/RacoonMediaServer/rms-packages/pkg/configuration"

// Configuration represents entire service configuration
type Configuration struct {
	Directory string
	Database  configuration.Database
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
