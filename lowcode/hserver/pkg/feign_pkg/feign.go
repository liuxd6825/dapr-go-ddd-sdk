package feign_pkg

import (
	"context"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/schema"
	"strings"
)

type Feign struct {
	server element.Server
	values map[string]any
}

type FeignParam struct {
	In       string         `json:"in"` // InParamType
	Required bool           `json:"required"`
	Type     string         `json:"type"`
	Desc     string         `json:"desc"`
	Schema   *schema.Schema `json:"schema"`
}

type FeignOptions struct {
	Method string //MethodType
	Url    string
	Params map[string]FeignParam
}

func New(server element.Server) *Feign {
	return newFeign(server)
}

func newFeign(server element.Server) *Feign {
	feign := &Feign{server: server}
	feign.values = map[string]any{
		"call":   feign.Call,
		"get":    feign.Get,
		"put":    feign.Put,
		"delete": feign.Delete,
		"patch":  feign.Patch,
		"post":   feign.Post,
	}
	return feign
}

func (f *Feign) NewProxy(vm *goja.Runtime, workPath string) *element.Proxy {
	return element.NewProxy(f.server, vm, workPath, f.values)
}

func (f *Feign) Call(ctx context.Context, options *FeignOptions, params map[string]any) any {
	return f.invoke(ctx, options, params)
}

func (f *Feign) Get(ctx context.Context, options *FeignOptions, params map[string]any) any {
	options.Method = "get"
	return f.invoke(ctx, options, params)
}

func (f *Feign) Post(ctx context.Context, options *FeignOptions, params map[string]any) any {
	options.Method = "post"
	return f.invoke(ctx, options, params)
}

func (f *Feign) Put(ctx context.Context, options *FeignOptions, params map[string]any) any {
	options.Method = "put"
	return f.invoke(ctx, options, params)
}

func (f *Feign) Delete(ctx context.Context, options *FeignOptions, params map[string]any) any {
	options.Method = "delete"
	return f.invoke(ctx, options, params)
}

func (f *Feign) Patch(ctx context.Context, options *FeignOptions, params map[string]any) any {
	options.Method = "patch"
	return f.invoke(ctx, options, params)
}

func (f *Feign) invoke(ctx context.Context, opts *FeignOptions, params map[string]any) any {
	var err error
	defer func() {
		if err = errors.GetRecoverError(err, recover()); err != nil {
			panic(err.Error())
		}
	}()

	if opts == nil {
		panic(errors.New("Feign.InvokeMethod() opts is nil"))
	}
	if opts.Url == "" {
		panic(errors.New("Feign.InvokeMethod() opts.URL is nil"))
	}
	if opts.Method == "" {
		panic(errors.New("Feign.InvokeMethod() opts.Method is nil"))
	}
	methodType := strings.ToLower(opts.Method)
	url, e := NewURLParser(opts.Url)
	if e != nil {
		return e
	}
	switch url.Protocol {
	case Protocol_Dapr:
		var request any
		var response any
		for k, v := range opts.Params {
			if v.In == InParamTypeBody.String() {
				request = params[k]
				break
			}
		}
		var data any
		data, err = dapr.GetDaprClient().InvokeService(ctx, url.ServiceName, url.Path, methodType, request, &response)
		if err != nil {
			panic(err)
		}
		return data
	default:
		err = errors.ErrorOf("unsupported protocol: %s", url.Protocol)
		panic(err.Error())
	}
}
