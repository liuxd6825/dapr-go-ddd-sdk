package service_test

import (
	_ "embed"
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/dao/impl/hugedao"
	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbevent"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
	"github.com/liuxd6825/jsonschema/v6"
)

const tenantId = "test"
const caseId = "1001"

//go:embed testfiles/company.json
var CompanySchema string

//go:embed testfiles/company-company.json
var CompanyCompanySchema string

var records = &Records{}

var companySch *jsonschema.Schema
var companyDBSch *dbschema.DBSchema
var companyCompanySch *jsonschema.Schema
var companyCompanyDBSch *dbschema.DBSchema

func getCompanySchema() *jsonschema.Schema {
	return schema.NewJsonSchemaWithJson("company.json", CompanySchema)
}

func getCompanyCompanySchema() *jsonschema.Schema {
	return schema.NewJsonSchemaWithJson("company_company.json", CompanyCompanySchema)
}

func newService() *service.MasterService {
	service := service.NewMasterService()

	companySch = getCompanySchema()
	companyDBSch = dbschema.NewDBSchemaWithJsonSchema(companySch)

	companyMetaExt := schema.GetMetaExtension(companySch)
	service.AddDao(companySch, companyDBSch, companyMetaExt, companyMetaExt.GetGraph())

	companyCompanySch = getCompanyCompanySchema()
	companyCompanyDBSch = dbschema.NewDBSchemaWithJsonSchema(companyCompanySch)

	companyCompanyMetaExt := schema.GetMetaExtension(companyCompanySch)
	service.AddDao(companyCompanySch, companyCompanyDBSch, companyCompanyMetaExt, companyCompanyMetaExt.GetGraph())

	return service
}

func Test_Case1(t *testing.T) {
	gp.Try(func() error {
		service := newService()
		//service.ClearAll(ctx)

		// 创建节点
		a := records.GetCreateCompany("a", "a", companyDBSch)
		service.Create(a)

		a_c := records.GetCreateCompanyCompany("a_c", "a", "子公司", "c", companyCompanyDBSch)
		service.Create(a_c)

		b := records.GetCreateCompany("b", "b", companyDBSch)
		service.Create(b)

		b_c := records.GetCreateCompanyCompany("b_c", "b", "子公司", "c", companyCompanyDBSch)
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

func Test_DeleteById(t *testing.T) {
	nodeDao := hugedao.NewMasterNodeDao([]string{"company_test"}, nil, nil)
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

type Records struct {
}

func (r *Records) GetCreateCompany(id string, name string, dbSch *dbschema.DBSchema) *dbevent.CDCRecord {
	record := &dbevent.CDCRecord{
		OpType: "c",
		DB:     "master",
		Table:  "company",
		After: map[string]any{
			"id":        id,
			"name":      name,
			"tenant_id": tenantId,
			"case_id":   caseId,
		},
		DBSchema: dbSch,
	}
	return record
}

func (r *Records) GetUpdateCompany(id, newName, oldName string) *dbevent.CDCRecord {
	record := &dbevent.CDCRecord{
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

func (r *Records) GetCreateCompanyCompany(id, startId, relType, name string, dbSch *dbschema.DBSchema) *dbevent.CDCRecord {
	record := &dbevent.CDCRecord{
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
		DBSchema: dbSch,
	}
	return record
}

func (r *Records) GetUpdateCompanyCompany(id, startId, newRelType, oldRelType, newName, oldName string) *dbevent.CDCRecord {
	record := &dbevent.CDCRecord{
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
