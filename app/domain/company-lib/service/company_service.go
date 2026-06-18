package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/company-lib/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/company-lib/model"
)

type CompanyService struct {
	dao *dao.CompanyDao
}

var companyService *CompanyService
var companyServiceOnce sync.Once

func NewCompanyService() *CompanyService {
	companyServiceOnce.Do(func() {
		companyService = &CompanyService{
			dao: dao.NewCompanyDao(),
		}
	})
	return companyService
}

func (s *CompanyService) Search(ctx context.Context, query *model.CompanyQuery) (*model.CompanyQueryResult, error) {
	dsl, err := buildCompanyQuery(query)
	if err != nil {
		return nil, err
	}

	result, err := s.dao.Search(ctx, dsl)
	if err != nil {
		return nil, err
	}

	return toCompanyResult(result), nil
}

func (s *CompanyService) FullTextSearch(ctx context.Context, query *model.FullTextSearchQuery) (*model.CompanyQueryResult, error) {
	dsl, err := buildFullCompanyQuery(query)
	if err != nil {
		return nil, err
	}

	result, err := s.dao.Search(ctx, dsl)
	if err != nil {
		return nil, err
	}

	return toCompanyResult(result), nil
}

func (s *CompanyService) GetById(ctx context.Context, id string) (*model.Company, error) {
	return s.dao.GetById(ctx, id)
}

func (s *CompanyService) Create(ctx context.Context, company *model.Company) error {
	return s.dao.Create(ctx, company)
}

func (s *CompanyService) Delete(ctx context.Context, id string) error {
	return s.dao.Delete(ctx, id)
}

func (s *CompanyService) BulkCreate(ctx context.Context, companies []*model.Company) (int, error) {
	return s.dao.BulkCreate(ctx, companies)
}

func (s *CompanyService) CreateIndex(ctx context.Context) error {
	return s.dao.CreateIndex(ctx)
}

func toCompanyResult(result *model.SearchResult) *model.CompanyQueryResult {
	companies := make([]*model.Company, 0, len(result.Hits))
	for _, hit := range result.Hits {
		company := &model.Company{}
		data, err := json.Marshal(hit)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(data, company); err != nil {
			continue
		}
		companies = append(companies, company)
	}

	return &model.CompanyQueryResult{
		Data:       companies,
		TotalCount: result.Total,
	}
}

func buildCompanyQuery(query *model.CompanyQuery) (string, error) {
	var clauses []string

	if query.Name != "" {
		clauses = append(clauses, buildMatchQuery("name", query.Name))
	}
	if query.RegNo != "" {
		clauses = append(clauses, buildMatchQuery("reg_no", query.RegNo))
	}
	if query.OperStatus != "" {
		clauses = append(clauses, buildMatchQuery("oper_status", query.OperStatus))
	}
	if query.CreditCode != "" {
		clauses = append(clauses, buildMatchQuery("credit_code", query.CreditCode))
	}
	if query.Iden != "" {
		clauses = append(clauses, buildMatchQuery("iden", query.Iden))
	}
	if query.ApprDate != "" {
		clauses = append(clauses, buildMatchQuery("appr_date", query.ApprDate))
	}
	if query.CreateDate != "" {
		clauses = append(clauses, buildMatchQuery("create_date", query.CreateDate))
	}
	if query.Qual != "" {
		clauses = append(clauses, buildMatchQuery("qual", query.Qual))
	}
	if query.EntType != "" {
		clauses = append(clauses, buildMatchQuery("ent_type", query.EntType))
	}
	if query.RegAuthority != "" {
		clauses = append(clauses, buildMatchQuery("reg_authority", query.RegAuthority))
	}
	if query.EngName != "" {
		clauses = append(clauses, buildMatchQuery("eng_name", query.EngName))
	}
	if query.Addr != "" {
		clauses = append(clauses, buildMatchQuery("addr", query.Addr))
	}
	if query.SpecificAddr != "" {
		clauses = append(clauses, buildMatchQuery("specific_addr", query.SpecificAddr))
	}
	if query.LegalPerson != "" {
		clauses = append(clauses, buildMatchQuery("legal_person", query.LegalPerson))
	}

	if len(clauses) == 0 {
		return `{"query": {"match_all": {}}}`, nil
	}

	return fmt.Sprintf(`{"query": {"bool": {"must": [%s]}}}`, strings.Join(clauses, ",")), nil
}

func buildFullCompanyQuery(query *model.FullTextSearchQuery) (string, error) {
	var clauses []string

	if query.Text != "" {
		clauses = append(clauses, buildMatchQuery("name", query.Text))
		clauses = append(clauses, buildMatchQuery("reg_no", query.Text))
		clauses = append(clauses, buildMatchQuery("oper_status", query.Text))
		clauses = append(clauses, buildMatchQuery("credit_code", query.Text))
		clauses = append(clauses, buildMatchQuery("iden", query.Text))
		clauses = append(clauses, buildMatchQuery("appr_date", query.Text))
		clauses = append(clauses, buildMatchQuery("create_date", query.Text))
		clauses = append(clauses, buildMatchQuery("qual", query.Text))
		clauses = append(clauses, buildMatchQuery("ent_type", query.Text))
		clauses = append(clauses, buildMatchQuery("reg_authority", query.Text))
		clauses = append(clauses, buildMatchQuery("eng_name", query.Text))
		clauses = append(clauses, buildMatchQuery("addr", query.Text))
		clauses = append(clauses, buildMatchQuery("specific_addr", query.Text))
		clauses = append(clauses, buildMatchQuery("legal_person", query.Text))
		clauses = append(clauses, buildWildcardQuery("bus_info", query.Text))
		clauses = append(clauses, buildWildcardQuery("sh_info", query.Text))
		clauses = append(clauses, buildWildcardQuery("key_person", query.Text))
	}

	if len(clauses) == 0 {
		return `{"query": {"match_all": {}}}`, nil
	}

	return fmt.Sprintf(`{"query": {"bool": {"minimum_should_match": 1,"should": [%s]}}}`, strings.Join(clauses, ",")), nil
}

func buildMatchQuery(field, value string) string {
	return fmt.Sprintf(`{"match": {"%s": "%s"}}`, field, escapeQueryValue(value))
}

func buildWildcardQuery(field, value string) string {
	return fmt.Sprintf(`{"wildcard": {"%s": "*%s*"}}`, field, escapeQueryValue(value))
}

func escapeQueryValue(value string) string {
	value = strings.ReplaceAll(value, `"`, `\"`)
	value = strings.ReplaceAll(value, `\`, `\\`)
	return value
}
