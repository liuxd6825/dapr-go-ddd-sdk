package feign_pkg

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"strings"
)

type Feign struct {
	server pkg.Server
}

func New(server pkg.Server) *Feign {
	return &Feign{server: server}
}

func (f *Feign) Call(ctx context.Context, options *Options, params map[string]any) *common.Result[any] {
	return f.InvokeMethod(ctx, options, params)
}

func (f *Feign) Get(ctx context.Context, options *Options, params map[string]any) *common.Result[any] {
	options.Method = "get"
	return f.InvokeMethod(ctx, options, params)
}

func (f *Feign) Post(ctx context.Context, options *Options, params map[string]any) *common.Result[any] {
	options.Method = "post"
	return f.InvokeMethod(ctx, options, params)
}

func (f *Feign) Put(ctx context.Context, options *Options, params map[string]any) *common.Result[any] {
	options.Method = "put"
	return f.InvokeMethod(ctx, options, params)
}

func (f *Feign) Delete(ctx context.Context, options *Options, params map[string]any) *common.Result[any] {
	options.Method = "delete"
	return f.InvokeMethod(ctx, options, params)
}

func (f *Feign) Patch(ctx context.Context, options *Options, params map[string]any) *common.Result[any] {
	options.Method = "patch"
	return f.InvokeMethod(ctx, options, params)
}

func (f *Feign) InvokeMethod(ctx context.Context, opts *Options, params map[string]any) (res *common.Result[any]) {
	var err error
	defer func() {
		if err = errors.GetRecoverError(err, recover()); err != nil {
			res = common.NewResult[any](nil, err)
		}
	}()

	if opts == nil {
		return common.NewResult[any](nil, errors.New("Feign.InvokeMethod() opts is nil"))
	}
	if opts.Url == "" {
		return common.NewResult[any](nil, errors.New("Feign.InvokeMethod() opts.URL is nil"))
	}
	if opts.Method == "" {
		return common.NewResult[any](nil, errors.New("Feign.InvokeMethod() opts.Method is nil"))
	}
	methodType := strings.ToLower(opts.Method)
	url, e := NewURLParser(opts.Url)
	if e != nil {
		return common.NewResult[any](nil, e)
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
		return common.NewResult[any](data, err)
	default:
		err = errors.ErrorOf("unsupported protocol: %s", url.Protocol)
		return common.NewResult[any](nil, err)
	}
}

type Protocol string

const (
	Protocol_Dapr  Protocol = "dapr"
	Protocol_Http  Protocol = "http"
	Protocol_Grpc  Protocol = "grpc"
	Protocol_Https Protocol = "https"
)

func (p Protocol) String() string {
	return string(p)
}

type URLParser struct {
	Protocol    Protocol `json:"protocol,omitempty"`
	ServiceName string   `json:"serviceName,omitempty"`
	Port        string   `json:"port,omitempty"`
	Path        string   `json:"path,omitempty"`
}

// NewURLParser
//
//	@Description:
//	@param str   http://user-service:8080/api/v1/users
//	@return *URLParser
func NewURLParser(str string) (*URLParser, error) {
	protocol := Protocol_Dapr
	s := strings.ToLower(str)
	if strings.HasPrefix(s, Protocol_Dapr.String()+"://") {
		protocol = Protocol_Dapr
	} else if strings.HasPrefix(s, Protocol_Http.String()+"://") {
		protocol = Protocol_Http
	} else if strings.HasPrefix(s, Protocol_Grpc.String()+"://") {
		protocol = Protocol_Grpc
	} else if strings.HasPrefix(s, Protocol_Https.String()+"://") {
		protocol = Protocol_Https
	}
	serviceName := ""
	port := ""
	path := ""
	s = str[len(protocol.String())+3:]
	list := strings.Split(s, "/")
	count := len(list)
	if count == 0 {
		return nil, errors.New("")
	} else if count >= 1 { // user-service:8080
		l := strings.Split(list[0], ":")
		c := len(l)
		if c == 1 {
			serviceName = l[0]
		} else if c >= 2 {
			serviceName = l[0]
			port = l[1]
		}
		path = "/" + strings.Join(list[1:], "/")
	}

	return &URLParser{Protocol: Protocol(protocol), ServiceName: serviceName, Port: port, Path: path}, nil
}
