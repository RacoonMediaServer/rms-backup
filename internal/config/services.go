package config

type Nextcloud struct {
	Container string
	User      string
	Database  string
	Data      string
}

type Postgres struct {
	Container string
}

type Services struct {
	Postgres  Postgres
	Nextcloud Nextcloud
}
