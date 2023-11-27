package restapp

import (
	"context"
	"fmt"
)

type DbManager interface {
	Create(ctx context.Context, table *Table, env *EnvConfig, options *CreateOptions)
	Update(ctx context.Context, table *Table, env *EnvConfig, options *UpdateOptions)
}

type InitOptions struct {
}

type CreateOptions struct {
	InitOptions
	Prefix string
}

type UpdateOptions struct {
	InitOptions
}

func NewCreateOptions() *CreateOptions {
	return &CreateOptions{}
}

func NewUpdateOptions() *UpdateOptions {
	return &UpdateOptions{}
}

func (o *CreateOptions) SetPrefix(v string) *CreateOptions {
	o.Prefix = v
	return o
}

func (o *InitOptions) Log(err error) {
	if err != nil {
		fmt.Printf(" error: %s ; \r\n", err.Error())
	}
}

func (o *InitOptions) Printf(format string, any2 ...any) {
	fmt.Printf(format, any2...)
}
