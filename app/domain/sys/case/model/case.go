package model

import "time"

type Case struct {
	BaseModel     `bson:",inline"`
	CaseTypeId    string     `json:"caseTypeId" gorm:"case_type_id" bson:"case_type_id"`
	Name          string     `json:"name" gorm:"name" bson:"name"`
	Code          string     `json:"code" gorm:"code" bson:"code"`
	FilingTime    *time.Time `json:"filingTime" gorm:"filing_time" bson:"filing_time"`
	SubjectName   string     `json:"subjectName" gorm:"subject_name" bson:"subject_name"`
	Investigators string     `json:"investigators" gorm:"investigators" bson:"investigators"`
	Status        string     `json:"status" gorm:"status"  bson:"status"`
}

func NewCase() (*Case, error) {
	return &Case{}, nil
}
