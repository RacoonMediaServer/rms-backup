package backup

import (
	"context"
	rms_backup "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-backup"
	"go-micro.dev/v4/logger"
	"sync"
	"sync/atomic"
	"time"
)

// Engine is an entity which can run Instruction
type Engine struct {
	l         logger.Logger
	isRunning atomic.Bool

	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	timeout time.Duration

	OnReady func(r Report)

	state state
}

func NewEngine() *Engine {
	return &Engine{
		l: logger.DefaultLogger.Fields(map[string]interface{}{"from": "backup"}),
	}
}

func (e *Engine) SetTimeout(timeout time.Duration) {
	e.timeout = timeout
}

func (e *Engine) Launch(ctx Context, backupType rms_backup.BackupType, instruction Instruction) bool {
	if !e.isRunning.CompareAndSwap(false, true) {
		return false
	}

	now := time.Now()
	e.state.setInProgress(backupType, genFileName(backupType, now), now)

	if e.timeout == 0 {
		e.ctx, e.cancel = context.WithCancel(context.Background())
	} else {
		e.ctx, e.cancel = context.WithTimeout(context.Background(), e.timeout)
	}
	ctx.ctx = e.ctx
	e.wg.Add(1)

	go func() {
		defer e.isRunning.Store(false)
		defer e.wg.Done()
		e.process(ctx, instruction)
	}()

	return true
}

func (e *Engine) GetReport() Report {
	return e.state.getReport()
}

func (e *Engine) Shutdown() {
	if !e.isRunning.Load() {
		return
	}
	e.cancel()
	e.wg.Wait()
}
