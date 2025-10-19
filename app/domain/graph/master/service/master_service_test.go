package service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/dao"
	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"testing"
)

var records = &Records{}

const tenantId = "test"
const caseId = "1001"

func Test_DeleteById(t *testing.T) {
	ctx := xtest.NewContext()
	envInv := xtest.NewEnvConfigNeo4j()
	env.SetEnv(envInv)
	nodeDao := dao.NewMasterNodeDao([]string{"company_test"}, nil)
	nodeDao.GetConfig().Env = envInv

	gp.Try(func() error {
		node := &model2.MasterNode{
			Id:       "XELWOSgMxzRASZBTGKDcSlHVGS",
			CaseId:   caseId,
			Name:     "星辰科技有限公司",
			Table:    "company",
			TenantId: tenantId,
		}
		nodeDao.Delete(ctx, node)
		return nil
	}).Catch(func(err error) {
		t.Error(err)
	})
}

func Test_Case1(t *testing.T) {
	gp.Try(func() error {
		envInv := xtest.NewEnvConfigNeo4j()
		env.SetEnv(envInv)

		ctx := xtest.NewContext()
		service := newService()
		service.clearAll(ctx)

		// 创建节点
		a := records.GetCreateCompany("a", "a")
		service.Create(a)

		a_c := records.GetCreateCompanyCompany("a_c", "a", "子公司", "c")
		service.Create(a_c)

		b := records.GetCreateCompany("b", "b")
		service.Create(b)

		b_c := records.GetCreateCompanyCompany("b_c", "b", "子公司", "c")
		service.Create(b_c)

		// 主数据改名
		a1 := records.GetUpdateCompany("a", "a1", "a")
		service.Update(a1)

		// 关系数据改名
		a_c1 := records.GetUpdateCompanyCompany("a_c", "a", "子公司", "子公司", "c1", "c")
		service.Update(a_c1)

		b_c1 := records.GetUpdateCompanyCompany("b_c", "b", "子公司", "子公司", "c1", "c")
		service.Update(b_c1)

		// 关系类型修改
		b_c2 := records.GetUpdateCompanyCompany("b_c", "b", "分公司", "子公司", "c1", "c1")
		service.Update(b_c2)

		return nil
	}).Catch(func(err error) {
		t.Error(err)
	})

}

func newService() *MasterService {
	env.SetEnv(xtest.NewEnvConfigNeo4j())

	companySch := xtest.GetCompanySchema()
	companyJsonSch := dbschema.NewDBSchemaWithJsonSchema(companySch)

	companyCompanySch := xtest.GetCompanyCompanySchema()
	companyCompanyJsonSch := dbschema.NewDBSchemaWithJsonSchema(companyCompanySch)

	service := NewMasterService()
	service.AddDao(companySch, companyJsonSch)
	service.AddDao(companyCompanySch, companyCompanyJsonSch)

	return service
}

type Records struct {
}

func (r *Records) GetCreateCompany(id string, name string) *model2.CDCRecord {
	record := &model2.CDCRecord{
		OpType: "c",
		DB:     "master",
		Table:  "company",
		After: map[string]any{
			"id":        id,
			"name":      name,
			"tenant_id": tenantId,
			"case_id":   caseId,
		},
	}
	return record
}

func (r *Records) GetUpdateCompany(id, newName, oldName string) *model2.CDCRecord {
	record := &model2.CDCRecord{
		OpType: "u",
		DB:     "master",
		Table:  "company",
		Before: map[string]any{
			"id":        id,
			"name":      oldName,
			"tenant_id": tenantId,
			"case_id":   caseId,
		},
		After: map[string]any{
			"id":        id,
			"name":      newName,
			"tenant_id": tenantId,
			"case_id":   caseId,
		},
	}
	return record
}

func (r *Records) GetCreateCompanyCompany(id, startId, relType, name string) *model2.CDCRecord {
	record := &model2.CDCRecord{
		OpType: "c",
		DB:     "master",
		Table:  "company_company",
		After: map[string]any{
			"id":            id,
			"company_id":    startId,
			"relation_type": relType,
			"name":          name,
			"tenant_id":     tenantId,
			"case_id":       caseId,
		},
	}
	return record
}

func (r *Records) GetUpdateCompanyCompany(id, startId, newRelType, oldRelType, newName, oldName string) *model2.CDCRecord {
	record := &model2.CDCRecord{
		OpType: "u",
		DB:     "master",
		Table:  "company_company",
		Before: map[string]any{
			"id":            id,
			"company_id":    startId,
			"relation_type": oldRelType,
			"name":          oldName,
			"tenant_id":     tenantId,
			"case_id":       caseId,
		},
		After: map[string]any{
			"id":            id,
			"company_id":    startId,
			"relation_type": newRelType,
			"name":          newName,
			"tenant_id":     tenantId,
			"case_id":       caseId,
		},
	}
	return record
}
