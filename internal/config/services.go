package config

type Nextcloud struct {
	Container string
	User      string
	Database  string
	Data      string
}

type Database struct {
	Container string
}

type Services struct {
	Database  Database
	Nextcloud Nextcloud
}
