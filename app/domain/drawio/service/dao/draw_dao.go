package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/drawio/service/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
)

type DrawDao struct {
	idao.Dao[*model.Draw]
}

func NewDrawDao() *DrawDao {
	d := dao.NewDao[*model.Draw](&dao.NewConfig{
		DBKey:    "mysql",
		DBSchema: dbschema.NewDBSchemaWithStruct("draw", &model.Draw{}, "draw"),
		Env:      env.GetEnv(),
		IsCache:  true,
	})
	return &DrawDao{
		Dao: d,
	}
}
