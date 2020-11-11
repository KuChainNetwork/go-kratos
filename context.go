package sdk

import (
	"context"

	"github.com/tendermint/tendermint/libs/log"
)

type Context struct {
	context.Context

	lcdURL string
	rpcURL string
	logger log.Logger
}

func NewCtx(ctx context.Context) Context {
	return Context{
		Context: ctx,
	}
}

func (c Context) WithUrls(lcdURL string, rpcURL string) Context {
	c.lcdURL = lcdURL
	c.rpcURL = rpcURL
	return c
}

func (c Context) WithLogger(logger log.Logger) Context {
	c.logger = logger
	return c
}

func (c Context) LcdURL() string {
	return c.lcdURL
}

func (c Context) RpcURL() string {
	return c.rpcURL
}

func (c Context) Logger() log.Logger {
	return c.logger
}
