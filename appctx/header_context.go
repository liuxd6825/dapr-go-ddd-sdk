package appctx

import (
	"context"
)

type ctxHeaderKey struct {
}

type Header map[string][]string

var headerKey = ctxHeaderKey{}

func NewHeaderContext(parent context.Context, metadata Header) context.Context {
	header := make(Header)
	for k, v := range metadata {
		header[k] = v
	}
	return context.WithValue(parent, headerKey, metadata)
}

func GetHeader(ctx context.Context) (Header, bool) {
	if ctx == nil {
		return nil, false
	}
	header := ctx.Value(headerKey)
	mapData, ok := header.(Header)
	return mapData, ok
}
