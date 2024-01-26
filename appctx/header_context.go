package appctx

import (
	"context"
)

type ctxHeaderKey struct {
}

type Header map[string][]string

var headerKey = ctxHeaderKey{}

func NewHeaderContext(pCtx context.Context, header map[string][]string) context.Context {
	data := Header{}
	for k, v := range header {
		data[k] = v
	}
	return context.WithValue(pCtx, headerKey, data)
}

func SetHeaderContext(parent context.Context, header map[string][]string) context.Context {
	data, ok := GetHeader(parent)
	if !ok {
		data = Header{}
	}
	for k, v := range header {
		data[k] = append(data[k], v...)
	}
	if ok {
		return parent
	}
	return context.WithValue(parent, headerKey, header)
}

func GetHeader(ctx context.Context) (Header, bool) {
	if ctx == nil {
		return nil, false
	}
	header := ctx.Value(headerKey)
	mapData, ok := header.(Header)
	return mapData, ok
}
