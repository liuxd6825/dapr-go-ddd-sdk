package restapp

import (
	"context"
)

type GORMManager struct {
}

func (m *GORMManager) Create(ctx context.Context, table *Table, env *EnvConfig, options *CreateOptions) {

}

func (m *GORMManager) Update(ctx context.Context, table *Table, env *EnvConfig, options *UpdateOptions) {

}

func (m *GORMManager) GetInitScript(ctx context.Context, dbKey string, table []*Table, env *EnvConfig, options *CreateOptions) string {
	return ""
}

func NewGormManager() DbManager {
	return &GORMManager{}
}

func NewGormScriptManager() DbScriptManager {
	return &GORMManager{}
}
