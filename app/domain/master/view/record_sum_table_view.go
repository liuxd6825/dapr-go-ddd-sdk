package view

import "time"

type RecordSumTableView struct {
	Id           any        `json:"id" bson:"_id"`
	TenantId     string     `json:"tenantId" bson:"tenant_id"`
	Name         string     `json:"name" bson:"name"`
	OppName      string     `json:"oppName" bson:"opp_name"`
	BeginDate    *time.Time `json:"beginDate" bson:"begin_date"`
	EndDate      *time.Time `json:"endDate" bson:"end_date"`
	Amount       float64    `json:"amount" bson:"amount"`
	Count        int        `json:"count" bson:"count"`
	OppAcctCount int        `json:"oppAcctCount" bson:"opp_acct_count"`
	Ccy          string     `json:"ccy" bson:"ccy"`
}

func NewRecordSumTableView() *RecordSumTableView {
	return &RecordSumTableView{}
}

func (v *RecordSumTableView) SetTenantId(val string) {
	v.TenantId = val
}

func (v *RecordSumTableView) GetTenantId() string {
	return v.TenantId
}
