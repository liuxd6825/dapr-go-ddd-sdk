package dao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/company-lib/elastic"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/company-lib/model"
)

type CompanyDao struct {
	client *elastic.ElasticClient
}

func NewCompanyDao() *CompanyDao {
	client, _ := elastic.NewElasticClient()
	return &CompanyDao{client: client}
}

const companyIndex = "company"

func (d *CompanyDao) Search(ctx context.Context, query string) (*model.SearchResult, error) {
	result, err := d.client.Search(ctx, companyIndex, query)
	if err != nil {
		return nil, err
	}
	return &model.SearchResult{
		Total: result.Total,
		Hits:  result.Hits,
	}, nil
}

func (d *CompanyDao) GetById(ctx context.Context, id string) (*model.Company, error) {
	var company model.Company
	err := d.client.Get(ctx, companyIndex, id, &company)
	if err != nil {
		return nil, err
	}
	return &company, nil
}

func (d *CompanyDao) Create(ctx context.Context, company *model.Company) error {
	doc := map[string]interface{}{
		"name":            company.Name,
		"reg_no":          company.RegNo,
		"oper_status":     company.OperStatus,
		"credit_code":     company.CreditCode,
		"iden":            company.Iden,
		"appr_date":       company.ApprDate,
		"create_date":     company.CreateDate,
		"qual":            company.Qual,
		"ent_type":        company.EntType,
		"reg_authority":   company.RegAuthority,
		"eng_name":        company.EngName,
		"addr":            company.Addr,
		"specific_addr":   company.SpecificAddr,
		"legal_person":    company.LegalPerson,
		"bus_info":        company.BusInfo,
		"sh_info":         company.ShInfo,
		"key_person":      company.KeyPerson,
		"is_transform":    company.IsTransform,
		"source_url":      company.SourceUrl,
		"source_id":       company.SourceId,
		"source_system":   company.SourceSystem,
		"remark":          company.Remark,
		"created_at":      company.CreatedAt,
		"updated_at":      company.UpdatedAt,
		"prev_updated_at": company.PrevUpdatedAt,
	}
	return d.client.Index(ctx, companyIndex, company.Id, doc)
}

func (d *CompanyDao) Delete(ctx context.Context, id string) error {
	return d.client.Delete(ctx, companyIndex, id)
}

func (d *CompanyDao) BulkCreate(ctx context.Context, companies []*model.Company) (int, error) {
	docs := make([]map[string]interface{}, 0, len(companies))
	for _, company := range companies {
		doc := map[string]interface{}{
			"_id":           company.Id,
			"name":          company.Name,
			"reg_no":        company.RegNo,
			"oper_status":   company.OperStatus,
			"credit_code":   company.CreditCode,
			"iden":          company.Iden,
			"appr_date":     company.ApprDate,
			"create_date":   company.CreateDate,
			"qual":          company.Qual,
			"ent_type":      company.EntType,
			"reg_authority": company.RegAuthority,
			"eng_name":      company.EngName,
			"addr":          company.Addr,
			"specific_addr": company.SpecificAddr,
			"legal_person":  company.LegalPerson,
			"bus_info":      company.BusInfo,
			"sh_info":       company.ShInfo,
			"key_person":    company.KeyPerson,
			"is_transform":  company.IsTransform,
			"source_url":    company.SourceUrl,
			"source_id":     company.SourceId,
			"source_system": company.SourceSystem,
			"remark":        company.Remark,
		}
		docs = append(docs, doc)
	}
	return d.client.Bulk(ctx, companyIndex, docs)
}

func (d *CompanyDao) CreateIndex(ctx context.Context) error {
	mappings := map[string]interface{}{
		"properties": map[string]interface{}{
			"name":            map[string]interface{}{"type": "text", "analyzer": "ik_max_word"},
			"reg_no":          map[string]interface{}{"type": "text", "analyzer": "ik_max_word"},
			"oper_status":     map[string]interface{}{"type": "keyword"},
			"credit_code":     map[string]interface{}{"type": "text", "analyzer": "ik_max_word"},
			"iden":            map[string]interface{}{"type": "text", "analyzer": "ik_max_word"},
			"appr_date":       map[string]interface{}{"type": "date"},
			"create_date":     map[string]interface{}{"type": "date"},
			"qual":            map[string]interface{}{"type": "text", "analyzer": "ik_max_word"},
			"ent_type":        map[string]interface{}{"type": "text", "analyzer": "ik_max_word"},
			"reg_authority":   map[string]interface{}{"type": "text", "analyzer": "ik_max_word"},
			"eng_name":        map[string]interface{}{"type": "text", "analyzer": "ik_max_word"},
			"addr":            map[string]interface{}{"type": "text", "analyzer": "ik_max_word"},
			"specific_addr":   map[string]interface{}{"type": "text", "analyzer": "ik_max_word"},
			"legal_person":    map[string]interface{}{"type": "text", "analyzer": "ik_max_word"},
			"bus_info":        map[string]interface{}{"type": "wildcard"},
			"sh_info":         map[string]interface{}{"type": "wildcard"},
			"key_person":      map[string]interface{}{"type": "wildcard"},
			"is_transform":    map[string]interface{}{"type": "boolean"},
			"source_url":      map[string]interface{}{"type": "keyword"},
			"source_id":       map[string]interface{}{"type": "keyword"},
			"source_system":   map[string]interface{}{"type": "keyword"},
			"remark":          map[string]interface{}{"type": "text"},
			"created_at":      map[string]interface{}{"type": "date"},
			"updated_at":      map[string]interface{}{"type": "date"},
			"prev_updated_at": map[string]interface{}{"type": "date"},
		},
	}
	return d.client.CreateIndex(ctx, companyIndex, mappings)
}
