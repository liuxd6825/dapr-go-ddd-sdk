package appctx

import (
	"context"
)

type ctxServerKey struct {
}

type Server interface {
	SetResponseHeader(key string, value string)
	URLParamDefault(name, def string) string
}

var serverKey = ctxServerKey{}

func NewServerContext(parent context.Context, server Server) context.Context {
	return context.WithValue(parent, serverKey, server)
}

func GetServer(ctx context.Context) (Server, bool) {
	if ctx == nil {
		return nil, false
	}

	header := ctx.Value(serverKey)
	if header == nil {
		return nil, false
	}
	serverContext, ok := header.(Server)
	return serverContext, ok
}
