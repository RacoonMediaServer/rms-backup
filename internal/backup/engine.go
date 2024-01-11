package backup

import (
	"context"
	"go-micro.dev/v4/logger"
	"sync"
	"sync/atomic"
)

// Engine is an entity which can run Instruction
type Engine struct {
	l         logger.Logger
	isRunning atomic.Bool

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	OnReady func(r Report)

	state state
}

func NewEngine() *Engine {
	return &Engine{
		l: logger.DefaultLogger.Fields(map[string]interface{}{"from": "backup"}),
	}
}

func (e *Engine) Launch(ctx Context, instruction Instruction) bool {
	if !e.isRunning.CompareAndSwap(false, true) {
		return false
	}

	e.state.setInProgress()

	e.ctx, e.cancel = context.WithCancel(context.Background())
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
