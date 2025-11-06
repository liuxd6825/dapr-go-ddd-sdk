package event

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/element"
)

type Config struct {
	tag *element.FuncTag
}

func NewConfig(tag *element.FuncTag) *Config {
	return &Config{
		tag: tag,
	}
}

func (c *Config) Name() string {
	return c.tag.Attr("name")
}

func (c *Config) AppId() string {
	return c.tag.Attr("app-id")
}

func (c *Config) Ver() string {
	return c.tag.AttrOr("ver", "v1.0")
}

func (c *Config) Pubsub() string {
	return c.tag.Attr("pubsub")
}
