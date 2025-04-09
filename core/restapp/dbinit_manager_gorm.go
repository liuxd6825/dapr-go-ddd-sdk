package restapp

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"strings"
)

type GORMManager struct {
}

func (m *GORMManager) Create(ctx context.Context, table *Table, env *env.Env, options *CreateOptions) {

}

func (m *GORMManager) Update(ctx context.Context, table *Table, env *env.Env, options *UpdateOptions) {

}

func (m *GORMManager) GetScript(ctx context.Context, dbKey string, table []*Table, env *env.Env, options *CreateOptions) (*strings.Builder, error) {
	return nil, nil
}

func NewGormManager() DbManager {
	return &GORMManager{}
}

func NewGormScriptManager() DbScriptManager {
	return &GORMManager{}
}
