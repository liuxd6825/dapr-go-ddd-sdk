package restapp

import (
	"context"
	"strings"
)

type GORMManager struct {
}

func (m *GORMManager) Create(ctx context.Context, table *Table, env *EnvConfig, options *CreateOptions) {

}

func (m *GORMManager) Update(ctx context.Context, table *Table, env *EnvConfig, options *UpdateOptions) {

}

func (m *GORMManager) GetScript(ctx context.Context, dbKey string, table []*Table, env *EnvConfig, options *CreateOptions) (*strings.Builder, error) {
	return nil, nil
}

func NewGormManager() DbManager {
	return &GORMManager{}
}

func NewGormScriptManager() DbScriptManager {
	return &GORMManager{}
}
