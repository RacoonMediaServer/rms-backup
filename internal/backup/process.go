package backup

import (
	"context"
	"errors"
	"fmt"
	"github.com/RacoonMediaServer/rms-backup/internal/config"
	"go-micro.dev/v4/logger"
	"path/filepath"
)

func (e *Engine) process(ctx Context, instruction Instruction) {
	defer func() {
		if r := recover(); r != nil {
			e.l.Logf(logger.ErrorLevel, "panic: %s", r)
			e.state.addError(fmt.Errorf("%+v", r))
			e.state.setFailed()
			e.onReady()
		}
	}()

	var passed []Command
	var artifacts []string
	success := 0
	failed := 0
	totalOp := instruction.Operations() + 1
	completeOp := 0

	e.l.Logf(logger.InfoLevel, "Start instruction set '%s'", instruction.Title)
	for i, stage := range instruction.Stages {
		l := e.l.Fields(map[string]interface{}{"stageNo": i + 1})
		l.Logf(logger.InfoLevel, "Run stage '%s'", stage.Title)
		if err := e.runStage(l, ctx, &stage, &passed); err != nil {
			e.state.addError(err)
			failed++
			l.Logf(logger.ErrorLevel, "Stage failed: %s", err)
			if errors.Is(err, context.Canceled) {
				break
			}
		} else {
			success++
			artifacts = append(artifacts, stage.Artifacts...)
		}
		completeOp++
		e.state.setProgress(completeOp, totalOp)
	}

	var size uint64
	if len(artifacts) != 0 {
		var err error
		size, err = e.compressArtifacts(ctx, artifacts)
		if err != nil {
			err = fmt.Errorf("compress artifacts failed: %w", err)
			e.l.Log(logger.ErrorLevel, err)
			e.state.addError(err)
			success = 0
		}
		completeOp++
		e.state.setProgress(completeOp, totalOp)
	}

	for i := len(passed) - 1; i >= 0; i-- {
		if err := safeCommandCleanup(ctx, passed[i]); err != nil {
			e.l.Logf(logger.WarnLevel, "Clean error: %s", err)
		}
	}

	if success != 0 {
		e.state.setReady(size)
	} else {
		e.state.setFailed()
	}

	e.onReady()

	e.l.Logf(logger.InfoLevel, "DONE. passed: %d, failed: %d", success, failed)
}

func (e *Engine) onReady() {
	if e.OnReady != nil {
		err := safeCall(func() error {
			e.OnReady(e.state.getReport())
			return nil
		})
		if err != nil {
			e.l.Logf(logger.ErrorLevel, "Call ready() callback failed: %s", err)
		}
	}
}

func (e *Engine) runStage(l logger.Logger, ctx Context, stage *Stage, passed *[]Command) error {
	var reason error

iterateStages:
	for i, cmd := range stage.Commands {
		l := l.Fields(map[string]interface{}{"cmd": i + 1})
		ctx.l = l
		l.Logf(logger.InfoLevel, "Run command '%s'", cmd.Title())

		select {
		case <-e.ctx.Done():
			l.Logf(logger.WarnLevel, "Canceled")
			reason = context.Canceled
			break iterateStages
		default:
		}

		if err := safeCommandExecute(ctx, cmd); err != nil {
			l.Logf(logger.ErrorLevel, "Command failed: %s", err)
			reason = err
			break
		}

		l.Log(logger.InfoLevel, "Done")
		*passed = append(*passed, cmd)
	}

	return reason
}

func (e *Engine) compressArtifacts(ctx Context, artifacts []string) (size uint64, err error) {
	ctx.Log().Logf(logger.InfoLevel, "Compress artifacts [ %d ]...", len(artifacts))
	archive := filepath.Join(config.Config().Directories.Backups, e.state.getReport().FileName)
	absArtifacts := make([]string, len(artifacts))
	for i := range artifacts {
		if !filepath.IsAbs(artifacts[i]) {
			absArtifacts[i] = filepath.Join(config.Config().Directories.Artifacts, artifacts[i])
		} else {
			absArtifacts[i] = artifacts[i]
		}
	}
	size, err = e.compressor.Compress(ctx, archive, artifacts)
	return
}
