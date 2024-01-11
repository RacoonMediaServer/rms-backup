package builder

import (
	"errors"
	"github.com/RacoonMediaServer/rms-backup/internal/backup"
	"time"
)

type debugCommand struct {
	haveErrors bool
}

func (d debugCommand) Title() string {
	return "Debug Wait"
}

func (d debugCommand) Execute(ctx backup.Context) error {
	<-time.After(5 * time.Second)
	if d.haveErrors {
		return errors.New("ERROR")
	}
	return nil
}

func (d debugCommand) Cleanup(ctx backup.Context) error {
	return nil
}

func createFullBackup() backup.Instruction {
	p := backup.Instruction{Title: "FullBackup"}

	s := backup.Stage{
		Title:     "Backup db",
		Commands:  []backup.Command{&debugCommand{}},
		Artifacts: []string{"backup.db"},
	}
	p.Add(s)

	s = backup.Stage{
		Title:     "Backup nc",
		Commands:  []backup.Command{&debugCommand{haveErrors: true}},
		Artifacts: []string{"backup.nextcloud"},
	}
	p.Add(s)

	return p
}
