package restapp

import (
	"context"
	"strings"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
)

type GORMManager struct {
}

func (m *GORMManager) Create(ctx context.Context, table *Table, env *envs.Env, options *CreateOptions) {

}

func (m *GORMManager) Update(ctx context.Context, table *Table, env *envs.Env, options *UpdateOptions) {

}

func (m *GORMManager) GetScript(ctx context.Context, dbKey string, table []*Table, env *envs.Env, options *CreateOptions) (*strings.Builder, error) {
	return nil, nil
}

func NewGormManager() DbManager {
	return &GORMManager{}
}

func NewGormScriptManager() DbScriptManager {
	return &GORMManager{}
}
