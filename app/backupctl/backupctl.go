package main

import (
	"context"
	"errors"
	"fmt"
	rms_backup "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-backup"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)
import "github.com/urfave/cli/v2"

func main() {
	var command string
	var backupFile string
	var fullBackup bool
	service := micro.NewService(
		micro.Name("rms-backup.client"),
		micro.Flags(
			&cli.StringFlag{
				Name:        "command",
				Usage:       "Must be one of: backup, list, remove",
				Required:    true,
				Destination: &command,
			},
			&cli.StringFlag{
				Name:        "backup",
				Usage:       "backup file",
				Required:    false,
				Destination: &backupFile,
			},
			&cli.BoolFlag{
				Name:        "full",
				Usage:       "create full backup",
				Required:    false,
				Destination: &fullBackup,
			},
		),
	)
	service.Init()

	client := rms_backup.NewRmsBackupService("rms-backup", service.Client())

	switch command {
	case "backup":
		if err := backup(client, fullBackup); err != nil {
			panic(err)
		}
	case "list":
		if err := list(client); err != nil {
			panic(err)
		}
	case "status":
		if err := status(client); err != nil {
			panic(err)
		}
	case "remove":
		if err := remove(client, backupFile); err != nil {
			panic(err)
		}
	case "apply":
		if err := apply(client); err != nil {
			panic(err)
		}
	default:
		panic("unknown command: " + command)
	}
}

func backup(cli rms_backup.RmsBackupService, full bool) error {

	req := rms_backup.LaunchBackupRequest{Type: rms_backup.BackupType_Partial}
	if full {
		req.Type = rms_backup.BackupType_Full
	}

	resp, err := cli.LaunchBackup(context.Background(), &req, client.WithRequestTimeout(40*time.Second))
	if err != nil {
		return err
	}

	if resp.AlreadyLaunch {
		return errors.New("already running")
	}
	return nil
}

func list(cli rms_backup.RmsBackupService) error {
	result, err := cli.GetBackups(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}
	for _, b := range result.Backups {
		fmt.Println(b.FileName, b.Type.String(), time.Unix(int64(b.Date), 0), b.Size)
	}
	return nil
}

func remove(cli rms_backup.RmsBackupService, id string) error {
	_, err := cli.RemoveBackup(context.Background(), &rms_backup.RemoveBackupRequest{FileName: id})
	return err
}

func status(cli rms_backup.RmsBackupService) error {
	resp, err := cli.GetBackupStatus(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}
	fmt.Println(resp.Status, resp.Progress, time.Unix(int64(resp.LastTime), 0))
	return nil
}

func apply(cli rms_backup.RmsBackupService) error {
	settings := rms_backup.BackupSettings{
		Enabled:  true,
		Type:     rms_backup.BackupType_Full,
		Period:   rms_backup.BackupSettings_EveryMonth,
		Day:      11,
		Hour:     20,
		Password: nil,
	}
	_, err := cli.SetBackupSettings(context.Background(), &settings)
	return err
}
