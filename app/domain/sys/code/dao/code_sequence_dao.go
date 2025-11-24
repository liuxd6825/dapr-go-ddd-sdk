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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CodeSequenceDao struct {
	idao.Dao[*model.CodeSequence]
}

func NewCodeSequenceDao(dbKey string) *CodeSequenceDao {
	tableName := "sys_code_sequence"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.CodeSequence{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.CodeSequence](newCfg)
	daoVal := &CodeSequenceDao{Dao: baseDao}
	return daoVal
}

// NextSeq 获取下一个序列号
// 如果当天记录不存在，会自动创建并返回 1
func (d *CodeSequenceDao) NextSeq(ctx context.Context, typeCode, dateStr string) (int64, error) {
	if d.Dao.GetDbType() == idao.DbType_MongoDB.String() {
		return d.mongoNextSeq(ctx, typeCode, dateStr)
	} else {

	}
	return 0, nil
}

// NextSeq 获取下一个序列号
// 如果当天记录不存在，会自动创建并返回 1
func (d *CodeSequenceDao) mongoNextSeq(ctx context.Context, typeCode, dateStr string) (int64, error) {
	mDao, ok := d.Dao.GetStore().(store_mongodb.IMongoDao[*model.CodeSequence])
	if !ok {
		return 0, fmt.Errorf("[code_sequence.dao] dao does not implement IMongoDao")
	}
	tenantId := appctx.GetTenantId2(ctx)
	filter := bson.M{
		"tenant_id": tenantId,
		"type_code": typeCode,
		"date_str":  dateStr,
	}

	// $inc: 原子递增
	// $set: 更新修改时间
	update := bson.M{
		"$inc": bson.M{"max_seq": 1},
		"$set": bson.M{"updated_at": time.Now()},
	}

	// Upsert: true (不存在则插入)
	// ReturnDocument: After (返回更新后的值，即 1, 2, 3...)
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result model.CodeSequence
	err := mDao.GetCollection(ctx).FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		return 0, err
	}

	return result.MaxSeq, nil
}
