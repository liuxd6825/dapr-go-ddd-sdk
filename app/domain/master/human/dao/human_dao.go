package dao

import (
	"context"
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

// HumanDao 人员数据访问对象
type HumanDao struct {
	idao.Dao[*model.Human]
}

func NewHumanDao(dbKey string) *HumanDao {
	tableName := "human"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Human{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.Human](newCfg)
	return &HumanDao{Dao: baseDao}
}

// FindPagingByCaseId 按 caseId 分页查询人员
func (d *HumanDao) FindPagingByCaseId(ctx context.Context, qry store.FindPagingQuery, caseId string, opts ...idao.CallOptions) store.FindPagingResult[*model.Human] {
	qry.SetMustFilter(fmt.Sprintf("case_id=='%s'", caseId))
	return d.FindPaging(ctx, qry, opts...)
}

// FindPagingByTagId 按 caseId + tagId 分页查询人员
func (d *HumanDao) FindPagingByTagId(ctx context.Context, qry store.FindPagingQuery, caseId, tagId string, opts ...idao.CallOptions) store.FindPagingResult[*model.Human] {
	qry.SetMustFilter(fmt.Sprintf("case_id=='%s' and tag_id=contains='%s'", caseId, tagId))
	qry.SetPageSize(store.PagingMaxPageSize)
	return d.FindPaging(ctx, qry, opts...)
}

// CountByCaseId 按 caseId 统计人员数量
func (d *HumanDao) CountByCaseId(ctx context.Context, caseId string, opts ...idao.CallOptions) (int64, error) {
	rsqlStr := fmt.Sprintf("case_id=='%s'", caseId)
	return d.CountByRSQL(ctx, rsqlStr, opts...)
}