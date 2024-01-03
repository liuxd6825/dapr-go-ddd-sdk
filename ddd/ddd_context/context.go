package ddd_context

import (
	"context"
)

type ctxMetadataKey struct {
}

type ctxServerKey struct {
}

type ServerContext interface {
	//SetResponseHeader(key string, value string)
	//URLParamDefault(name, def string) string
}

func NewContext(parent context.Context, metadata map[string]string, serverCtx ServerContext) context.Context {
	ctx, _ := context.WithCancel(parent)
	ctx = setMetadata(ctx, metadata)
	return context.WithValue(ctx, ctxServerKey{}, serverCtx)
}

func GetMetadataContext(ctx context.Context) map[string]string {
	header := ctx.Value(ctxMetadataKey{})
	mapData, ok := header.(map[string]string)
	if ok {
		return mapData
	}
	mapData = make(map[string]string)
	return mapData
}

func GetServerContext(ctx context.Context) ServerContext {
	header := ctx.Value(ctxServerKey{})
	serverContext, ok := header.(ServerContext)
	if ok {
		return serverContext
	}
	return nil
}

/*
func SetResponseHeader(ctx context.Context, name string, value string) {
	GetServerContext(ctx).SetResponseHeader(name, value)
}

func SetMetadataValue(ctx context.Context, name string, value string) {
	header := ctx.Value(ctxMetadataKey{})
	mapData, ok := header.(map[string]string)
	if ok {
		mapData[name] = value
	}
}
*/

func setMetadata(ctx context.Context, metadata map[string]string) context.Context {
	return context.WithValue(ctx, ctxMetadataKey{}, metadata)
}
