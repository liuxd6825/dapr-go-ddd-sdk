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
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"
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

func (d *CodeTypeDao) GetByCode(ctx context.Context, code string, style model.CodeStyle) (*model.CodeType, error) {
	if d.GetDbType() == idao.DbType_MongoDB.String() {
		return d.mongoGetByCode(ctx, code, style)
	}
	return nil, errors.New("not found code")
}

func (d *CodeTypeDao) mongoGetByCode(ctx context.Context, code string, style model.CodeStyle) (*model.CodeType, error) {
	mDao, ok := d.GetStore().(store_mongodb.IMongoDao[*model.CodeType])
	if !ok {
		return nil, errors.New("")
	}
	tenantId := appctx.GetTenantId2(ctx)
	filter := bson.M{
		"tenant_id": tenantId,
		"code":      code,
	}
	dateFormat := style.DateFormat()

	id := idutils.NewId()
	update := bson.M{
		"$setOnInsert": bson.M{
			"id":          id,
			"tenant_id":   tenantId,
			"code":        code,
			"name":        code,       // 默认名称等于Code
			"prefix":      code + "-", // 默认前缀：比如 "MH" -> "MH-"
			"date_format": dateFormat, // 默认格式：YYMMDD
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
