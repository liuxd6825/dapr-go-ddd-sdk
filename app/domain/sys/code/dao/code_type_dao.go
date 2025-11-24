package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/code/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/store_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CodeTypeDao struct {
	idao.Dao[*model.CodeType]
}

func NewCodeTypeDao(dbKey string) *CodeTypeDao {
	tableName := "sys_code_type"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.CodeType{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.CodeType](newCfg)
	daoVal := &CodeTypeDao{Dao: baseDao}
	return daoVal
}

func (d *CodeTypeDao) GetByCode(ctx context.Context, code string) (*model.CodeType, error) {
	if d.GetDbType() == idao.DbType_MongoDB.String() {
		return d.mongoGetByCode(ctx, code)
	}
	return nil, errors.New("not found code")
}

func (d *CodeTypeDao) mongoGetByCode(ctx context.Context, code string) (*model.CodeType, error) {
	mDao, ok := d.GetStore().(store_mongodb.IMongoDao[*model.CodeType])
	if !ok {
		return nil, errors.New("")
	}
	tenantId := appctx.GetTenantId2(ctx)
	filter := bson.M{
		"tenant_id": tenantId,
		"code":      code,
	}

	// 定义默认值：只有在文档是新插入(Insert)的时候，这些字段才会被设置
	// 如果文档已存在，这些字段会被忽略，保证不会覆盖用户修改过的配置
	update := bson.M{
		"$setOnInsert": bson.M{
			"tenant_id":   tenantId,
			"code":        code,
			"name":        code,       // 默认名称等于Code
			"prefix":      code + "-", // 默认前缀：比如 "MH" -> "MH-"
			"date_format": "060102",   // 默认格式：YYMMDD
			"seq_length":  6,          // 默认长度：6位
			"created_at":  time.Now(),
			"updated_at":  time.Now(),
		},
	}

	// 设置 Upsert=true (不存在则插入)
	// 设置 ReturnDocument=After (返回处理后的文档，即刚创建的或已存在的)
	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var result model.CodeType
	// 执行原子操作
	err := mDao.GetCollection(ctx).FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create code type: %v", err)
	}

	return &result, nil

}
