package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/service/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
)

type DrawDao struct {
	idao.Dao[*model.Draw]
}

func NewDrawDao(dbKey string) *DrawDao {
	base := dao.NewDao[*model.Draw](&dao.DaoConfig{
		DBKey:    dbKey,
		DBSchema: dbschema.NewDBSchemaWithStruct("draw", &model.Draw{}, "draw"),
		Env:      env.GetEnv(),
		IsCache:  true,
	})
	return &DrawDao{
		Dao: base,
	}
}
