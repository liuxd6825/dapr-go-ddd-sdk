package restapp

import (
	"fmt"
	"strings"
)

type ItemError struct {
	name string
	msg  string
}

type ConfigError struct {
	items []ItemError
}

type EnvTypeError struct {
	msg string
}

func NewEnvTypeError(msg string) *EnvTypeError {
	return &EnvTypeError{
		msg: msg,
	}
}

func (c *ConfigError) Error() string {
	sb := strings.Builder{}
	for _, item := range c.items {
		sb.WriteString(fmt.Sprintf("%s %s", item.name, item.msg))
	}
	return sb.String()
}

func (c *ConfigError) ItemCount() int {
	return len(c.items)
}

func (e *EnvTypeError) Error() string {
	return e.msg
}
