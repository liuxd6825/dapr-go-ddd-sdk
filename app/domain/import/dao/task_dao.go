package dao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/field"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors/assert"
)

type TaskDao struct {
	idao.Dao[*model.Task]
}

func NewTaskDao(dbKey string) *TaskDao {
	tableName := "import_task"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Task{}, tableName)
	newCfg := &dao.NewConfig{
		DBKey:      dbKey,
		IsPubEvent: dao.IsFalse(),
		TableName:  tableName,
		DBSchema:   dbSch,
	}
	baseDao := dao.NewDao[*model.Task](newCfg)
	daoVal := &TaskDao{Dao: baseDao}
	return daoVal
}

func (r *TaskDao) UpdateLock(ctx context.Context, taskId string, lock bool, opts ...idao.CallOptions) error {
	rsqlBuilder := rsql.NewBuilder().And(rsql.Eq("id", taskId), rsql.Eq("lock", lock))
	data := map[string]any{
		"lock": lock,
	}
	return r.UpdateMapByRSQL(ctx, rsqlBuilder.Build(), data).GetError()
}

func (r *TaskDao) UpdateProgress(ctx context.Context, setFields *field.TaskUpdateProgressFields) error {
	if err := assert.NotEmpty(setFields.Id); err != nil {
		return err
	}
	data := map[string]any{
		"end_time": setFields.EndTime,
		"message":  setFields.Message,
	}
	if setFields.State != "" {
		data["state"] = setFields.State
	}
	if setFields.StartTime != nil {
		data["start_time"] = setFields.EndTime
	}
	/*inc := bson.D{
		{"complete_rows", setFields.CompleteRows},
	}*/
	//data := map[string]any{ set}, {"$inc", inc}}
	return r.Dao.UpdateMap(ctx, setFields.Id, data).GetError()
}

func (r *TaskDao) SetState(ctx context.Context, taskId string, state model.TaskState, message string) error {
	if err := assert.NotEmpty(taskId); err != nil {
		return err
	}
	rsqlBuilder := rsql.NewBuilder(func(b *rsql.Builder) rsql.Condition {
		return rsql.Eq("id", taskId)
	})
	data := map[string]any{
		"state":   state,
		"message": message,
	}
	return r.Dao.UpdateMapByRSQL(ctx, rsqlBuilder.Build(), data).GetError()
}
