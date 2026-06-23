package view

type RecordSumChartView struct {
	Id          any     `json:"id" bson:"_id"`
	TenantId    string  `json:"tenantId" bson:"tenant_id"`
	Amount      float64 `json:"amount" bson:"amount"`
	Count       int     `json:"count" bson:"count"`
	OppNum      int     `json:"oppNum" bson:"opp_num"`
	Year        int     `json:"year" bson:"year"`
	Month       int     `json:"month" bson:"month"`
	Day         int     `json:"day" bson:"day"`
	SummaryType string  `json:"summaryType" bson:"summary_type"` //year、month、day
	Ccy         string  `json:"ccy" bson:"ccy"`
}

func NewRecordSumChartView() *RecordSumChartView {
	return &RecordSumChartView{}
}

//func (v *RecordSumChartView) SetId(val string) {
//	v.Id = val
//}
//
//func (v *RecordSumChartView) GetId() string {
//	return v.Id
//}

func (v *RecordSumChartView) SetTenantId(val string) {
	v.TenantId = val
}

func (v *RecordSumChartView) GetTenantId() string {
	return v.TenantId
}
