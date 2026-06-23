package view

import (
	"fmt"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/model"
)

// RecordDayView
// @Description: 按年月日+MasterType+MasterId+开户名 汇总
type RecordDayView struct {
	Id       string `json:"id"  bson:"_id" validate:"required" desc:"行Id"` // 行Id
	TenantId string `json:"tenantId" bson:"tenant_id"`
	CaseId   string `json:"caseId" bson:"case_id"`
	GraphId  string `json:"graphId" bson:"graph_id"`

	MasterId   string `json:"masterId" bson:"master_id"`
	MasterType string `json:"masterType" bson:"master_type"`

	Name string `json:"name"  bson:"name" validate:"-" desc:"名称"` // 名称
	Acct string `json:"acct" bson:"acct" validate:"-" desc:"账号"`  // 账号

	OppName string `json:"oppName" bson:"opp_name" validate:"-" desc:"对方名称"`  // 对方名称
	OppAcct string `json:"oppAcct"  bson:"opp_acct" validate:"-" desc:"对方账号"` // 对方账号

	Date  time.Time `json:"date" time_format:"2006-01-02" bson:"date" desc:"交易日期"`
	Year  int       `json:"year" bson:"year" desc:"交易年"`                // 交易年
	Month int       `json:"month" bson:"month" desc:"交易月"`              // 交易月
	Day   int       `json:"day" bson:"day" desc:"交易日"`                  // 交易日
	Ccy   string    `json:"ccy"  bson:"ccy"  validate:"-" desc:"交易币种" ` // 交易币种

	Count  int64   `json:"count" bson:"count" desc:"数量"`
	Payout float64 `json:"payout"   bson:"payout"  validate:"-" desc:"借方发生额"` // 借方发生额（支取）
	Income float64 `json:"income"  bson:"income"   validate:"-" desc:"贷方发生额"` // 贷方发生额（收入）
	Amount float64 `json:"amount"   bson:"amount"  validate:"-" desc:"交易金额"`  // 交易金额
}

func NewRecordDayViewId(v *RecordDayView) string {
	return fmt.Sprintf("%s-%s-%s-%s-%d-%d-%d-%s", v.TenantId, v.CaseId, v.MasterType, v.MasterId, v.Year, v.Month, v.Day, v.Name)
}

func (v *RecordDayView) AutoSetId() {
	v.Id = NewRecordDayViewId(v)
}

func (v *RecordDayView) SetId(val string) {
	v.Id = val
}

func (v *RecordDayView) GetId() string {
	return v.Id
}

func (v *RecordDayView) SetTenantId(val string) {
	v.TenantId = val
}

func (v *RecordDayView) GetTenantId() string {
	return v.TenantId
}

// NewRecordDayView
// @Description: 创建一个双向的日期流水，以方便统计查询
// @param record
// @return *RecordDayView
// @return *RecordDayView
// @return error
func NewRecordDayView(record *model.Record) (*RecordDayView, error) {
	year, month, day := record.Date.Date()
	// 重新构造时间：时分秒纳秒设为 0，地点(Location)保持不变
	date := time.Date(year, month, day, 0, 0, 0, 0, record.Date.Location())
	d1 := &RecordDayView{
		TenantId: record.TenantId,
		CaseId:   record.CaseId,

		MasterType: record.MasterType,
		MasterId:   record.MasterId,
		//GraphId:  record.GraphId,

		Acct:    record.Acct,
		Name:    record.Name,
		OppAcct: record.OppAcct,
		OppName: record.OppName,

		Date:  date,
		Year:  year,
		Month: int(month),
		Day:   day,

		Payout: record.Payout,
		Income: record.Income,
		Amount: record.Amount,
		Ccy:    record.Ccy,
		Count:  1,
	}
	/*
		d2 := &RecordDayView{
			TenantId: record.TenantId,
			CaseId:   record.CaseId,

			MasterType: record.MasterType,
			MasterId:   record.MasterId,
			GraphId:    record.GraphId,

			Acct:    record.OppAcct,
			Name:    record.OppName,
			OppAcct: record.Acct,
			OppName: record.Name,

			Date:  date,
			Year:  year,
			Month: int(month),
			Day:   day,

			Payout: record.Income,
			Income: record.Payout,
			Amount: record.Amount,
			Ccy:    record.Ccy,
			Count:  1,
		}
	*/
	d1.AutoSetId()
	//	d2.AutoSetId()

	return d1, nil
}
