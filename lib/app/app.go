package app

import (
	"context"
	"sync"
)

var (
	_AppCtx *AppContext
	_Once   sync.Once
)

type AppContext struct {
	ctx    context.Context
	logger *logger
}

func (a *AppContext) Log() *logger {
	return a.logger
}

func (a *AppContext) Ctx() context.Context {
	return a.ctx
}

func Init(ctx context.Context, debugEnabled bool) {
	_Once.Do(func() {
		_AppCtx = &AppContext{
			ctx:    ctx,
			logger: newLogger(debugEnabled),
		}
	})
}

func Get() *AppContext {
	return _AppCtx
}
