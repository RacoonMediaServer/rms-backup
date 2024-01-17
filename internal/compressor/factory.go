package compressor

import "context"

type Compressor interface {
	Compress(ctx context.Context, targetArchive string, files []string) (size uint64, err error)
	Extension() string
}

type Format int

type Settings struct {
	Password *string
}

type SettingsProvider func() Settings

const (
	Format_7z Format = iota
)

func New(format Format, provider SettingsProvider) Compressor {
	switch format {
	case Format_7z:
		return &zip7{provider: provider}
	default:
		panic("Unimplemented format type")
	}
}
