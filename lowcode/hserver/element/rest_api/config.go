package rest_api

import "github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"

type RestConfig struct {
	tag *element.FuncTag
}

func NewRestConfig(tag *element.FuncTag) *RestConfig {
	cfg := RestConfig{
		tag: tag,
	}
	return &cfg
}

func (r *RestConfig) Tag() *element.FuncTag {
	return r.tag
}

func (r *RestConfig) Method() string {
	return r.tag.Attr("method")
}

func (r *RestConfig) Url() string {
	return r.tag.Attr("url")
}

func (r *RestConfig) AbsUrl() string {
	return r.tag.Attr("abs-url")
}

func (r *RestConfig) ParamsType() string {
	return r.tag.Attr("params-type")
}

func (r *RestConfig) ParamsUrl() string {
	return r.tag.Attr("params-url")
}
