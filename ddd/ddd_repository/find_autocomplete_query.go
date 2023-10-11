package ddd_repository

import "fmt"

type FindAutoCompleteQuery struct {
	FindPagingQuery
	CaseId         string `json:"caseId" bson:"case_id" gorm:"case_id"`
	Field          string `json:"field" bson:"field" gorm:"field"`
	Value          string `json:"value" bson:"value" gorm:"value"`
	AllowTotalRows bool   `json:"allowTotalRows"`
}

func (f *FindAutoCompleteQuery) GetMustWhere() string {
	if len(f.CaseId) > 0 && len(f.Field) > 0 {
		return fmt.Sprintf(`caseId=="%s" and %s~="%s"`, f.CaseId, f.Field, f.Value)
	}
	if len(f.Field) > 0 {
		return fmt.Sprintf(`%s~="%s"`, f.Field, f.Value)
	}
	return ""
}
