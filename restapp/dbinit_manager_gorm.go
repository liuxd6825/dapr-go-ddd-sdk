package restapp

import (
	"context"
)

type GORMManager struct {
}

func (m *GORMManager) Create(ctx context.Context, table *Table, env *EnvConfig, options *CreateOptions) {
	//TODO implement me
	panic("implement me")
}

func (m *GORMManager) Update(ctx context.Context, table *Table, env *EnvConfig, options *UpdateOptions) {
	//TODO implement me
	panic("implement me")
}

func NewGormManager() DbManager {
	return &GORMManager{}
}
