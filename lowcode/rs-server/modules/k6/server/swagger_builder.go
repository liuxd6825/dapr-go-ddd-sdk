package server

import (
	"fmt"
	"github.com/dop251/goja"
	swagger3 "github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/swagger/v3"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"strings"
)

func Swagger(args ...any) {
	fmt.Print(args...)
}

func Decorate(args ...any) {
	fmt.Print(args...)
}

type SwaggerBuilder struct {
}

func NewSwaggerBuilder() *SwaggerBuilder {
	return &SwaggerBuilder{}
}

func (s *SwaggerBuilder) Build(server *Server) (*swagger3.Swagger, error) {
	swagger := swagger3.NewSwagger()
	for _, service := range server.services {
		if err := s.addPaths(swagger, service); err != nil {
			return nil, err
		}
	}
	return swagger, nil
}

func (s *SwaggerBuilder) addPaths(swagger *swagger3.Swagger, service *Service) error {
	apis := service.Service.Get("$apis")
	if apis == nil {
		return nil
	}
	if obj, ok := apis.(*goja.Object); ok {
		for _, key := range obj.Keys() {
			if v, ok := obj.Get(key).(*goja.Object); ok {
				method, _, err := s.newPathMethod(service.Name, v)
				if err != nil {
					return err
				}
				path, ok := swagger.Paths[method.Path]
				if !ok {
					path = swagger3.Path{}
					swagger.Paths[method.Path] = path
				}
				path[method.Method] = method
			}
		}
	}
	return nil
}

func (s *SwaggerBuilder) getString(v *goja.Object, key string) string {
	return v.Get(key).String()
}

func (s *SwaggerBuilder) getHandleOptions(value goja.Value) (*HandleOptions, error) {
	if o, ok := value.(*goja.Object); ok {
		bytes, err := o.MarshalJSON()
		if err != nil {
			return nil, err
		}
		opts := &HandleOptions{}
		err = jsonutils.Unmarshal(bytes, opts)
		return opts, err
	}
	return nil, fmt.Errorf("value is not a HandleOptions object")
}

func (s *SwaggerBuilder) newPathMethod(serviceName string, obj *goja.Object) (*swagger3.PathMethod, *HandleOptions, error) {
	h, err := s.getHandleOptions(obj)
	if err != nil {
		return nil, nil, err
	}
	method := &swagger3.PathMethod{}
	method.OperationId = h.HandleName
	method.Summary = h.Description
	method.Resource = serviceName
	method.Path = h.Path
	method.Method = strings.ToLower(h.Method.String())
	for name, param := range h.Params {
		par, err := s.newParameter(&param, name)
		if err != nil {
			return nil, nil, err
		}
		if par.Schema == nil {
			par.Schema = &swagger3.Schema{Type: "string"}
		}
		par.Schema = param.Schema
		method.Parameters = append(method.Parameters, par)
	}
	return method, h, nil
}

func (s *SwaggerBuilder) newParameter(param *RequestParam, name string) (*swagger3.Parameter, error) {
	par := &swagger3.Parameter{
		Name:        name,
		In:          param.In,
		Description: param.Description,
		Required:    param.Required,
		Example:     param.Example,
		Schema:      param.Schema,
	}
	return par, nil
}
