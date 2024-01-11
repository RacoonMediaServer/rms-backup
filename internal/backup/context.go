package backup

import (
	"context"
	"go-micro.dev/v4/logger"
	"time"
)

type Context struct {
	l   logger.Logger
	ctx context.Context
}

func NewContext() Context {
	return Context{
		l:   logger.DefaultLogger,
		ctx: context.Background(),
	}
}

func (c Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c Context) Deadline() (deadline time.Time, ok bool) {
	deadline, ok = c.ctx.Deadline()
	return
}

func (c Context) Err() error {
	return c.ctx.Err()
}

func (c Context) Value(key any) any {
	return c.ctx.Value(key)
}

func (c Context) Log() logger.Logger {
	return c.l
}
