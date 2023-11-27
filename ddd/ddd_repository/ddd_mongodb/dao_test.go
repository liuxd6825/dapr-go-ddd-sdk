package ddd_mongodb

import (
	"context"
	"encoding/json"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

func TestMapper_Search(t *testing.T) {
	ctx := context.Background()
	mdb, coll := newCollection("record_ie")
	mapper := NewDao[*Record](func() (mongodb *MongoDB, collection *mongo.Collection) {
		return mdb, coll
	})

	qry := ddd_repository.NewFindPagingQuery()
	qry.SetFilter("name=='梁瑞梅' and batch_id=='bd82aaf8-1654-4680-ba29-4e16eb66c29f'")
	qry.SetTenantId("test")
	qry.SetPageSize(20)
	qry.SetGroupCols(ddd_repository.NewGroupCols().Add("name", types.DataTypeString).Add("oppName", types.DataTypeString).Cols())
	qry.SetValueCols(ddd_repository.NewValueCols().Add("amount", ddd_repository.AggFuncSum).Cols())
	qry.SetGroupKeys([]any{"梁瑞梅"})

	res := mapper.FindPaging(ctx, qry)
	if res.Error == nil {
		logObject(t, "Data: ", res.Data)
		logObject(t, "Sum: ", res.SumData)
	} else {
		t.Error(res.Error)
	}

}

func TestDao_CreateIndexes(t *testing.T) {
	ctx := context.Background()
	mdb, coll := newCollection("test_create_index")
	mapper := NewDao[*Index](func() (mongodb *MongoDB, collection *mongo.Collection) {
		return mdb, coll
	})

	err := mapper.CreateIndexes(ctx)
	if err != nil {
		t.Error(err)
	}
}

func logObject(t *testing.T, label string, obj any) {
	jsonText, err := json.Marshal(obj)
	if err != nil {
		t.Error(label, err)
	}
	t.Log(label, string(jsonText))
}

type Index struct {
	Id        string `bson:"_id" `
	TenantId  string
	Name      string `bson:"name" index:"" `
	Asc       int64  `bson:"asc" index:" asc"`
	Desc      int64  `bson:"desc" index:" desc "`
	Unique    string `index:"unique"`
	AscUnique string `bson:"asc_unique" index:"asc, unique "`
}

type Record struct {
	Id     string `json:"id,omitempty"  bson:"_id"  index:""  validate:"required"  description:"id"` // 行Id
	RowNum int64  `json:"rowNum,omitempty" bson:"rowNum" description:"租户标识"`

	TenantId string   `json:"tenantId,omitempty" bson:"tenant_id" description:"租户标识"`
	DocId    string   `json:"docId,omitempty" bson:"doc_id"  description:"文档id"`
	FileId   string   `json:"fileId,omitempty" bson:"file_id"  description:"文件id"`
	BatchId  string   `json:"batchId,omitempty" bson:"batch_id"  description:"批id"`
	CaseId   string   `json:"caseId,omitempty" bson:"case_id" description:"案件id"`
	Iden     string   `json:"iden,omitempty"  bson:"iden"   validate:"-" description:"我方标识"`            // 标识
	Name     string   `json:"name,omitempty"   bson:"name"  validate:"-" description:"我方名称"`            // 名称
	Acct     string   `json:"acct,omitempty"   bson:"acct"  validate:"-" description:"我方账号"`            // 账号
	AcctType string   `json:"acctType,omitempty"   bson:"acct_type"  validate:"-" description:"我方账号类型"` // 账号类型
	Category string   `json:"category,omitempty"   bson:"category"  validate:"-" description:"我方类别"`    // 类别Id 公司或个人
	BankName string   `json:"bankName"  bson:"bank_name"   validate:"-" description:"我方开户银行"`           // 开户银行
	Balance  *float64 `json:"balance,omitempty"   bson:"balance"  validate:"-" description:"我方余额账户"`    // 余额账户

	OppIden     string `json:"oppIden,omitempty"   bson:"opp_iden"  validate:"-" description:"对方标识"`            // 对方标识
	OppName     string `json:"oppName,omitempty"   bson:"opp_name"  validate:"-" description:"对方名称"`            // 对方名称
	OppAcct     string `json:"oppAcct,omitempty"   bson:"opp_acct"  validate:"-" description:"对方账号"`            // 对方账号
	OppAcctType string `json:"oppAcctType,omitempty"   bson:"opp_acct_type"  validate:"-" description:"对方账号类型"` // 对方账号类型
	OppCategory string `json:"oppCategory,omitempty"   bson:"opp_category"  validate:"-" description:"对方类别"`    // 对方类别
	OppBankName string `json:"oppBankName,omitempty"  bson:"opp_bank_name"   validate:"-" description:"对方开户银行"` // 对方开户银行

	Serial  string     `json:"serial,omitempty"   bson:"serial"  validate:"-" description:"流水号"`   // 流水号
	Payout  *float64   `json:"payout,omitempty"   bson:"payout"  validate:"-" description:"借方发生额"` // 借方发生额（支出）
	Income  *float64   `json:"income,omitempty"  bson:"income"   validate:"-" description:"贷方发生额"` // 贷方发生额（收入）
	Amount  *float64   `json:"amount,omitempty"   bson:"amount"  validate:"-" description:"交易金额"`  // 交易金额
	Date    *time.Time `json:"date,omitempty"   bson:"date"  validate:"-" description:"交易时间"`      // 交易时间
	Type    string     `json:"type,omitempty"   bson:"type"  validate:"-" description:"交易类型"`      // 交易类型
	Ccy     string     `json:"ccy,omitempty"  bson:"ccy"  validate:"-" description:"交易币种" `        // 交易币种
	Place   string     `json:"place,omitempty"   bson:"place"  validate:"-" description:"地点"`      // 交易地点
	Summary string     `json:"summary,omitempty"   bson:"summary"  validate:"-" description:"摘要"`  // 摘要
	Notes   string     `json:"notes,omitempty"   bson:"notes"  validate:"-" description:"备注"`      // 备注

	HasError bool `json:"hasError" bson:"hasError,omitempty"`
}

func (u *Record) GetTenantId() string {
	return u.TenantId
}

func (u *Record) SetTenantId(v string) {
	u.TenantId = v
}

func (u *Record) GetId() string {
	return string(u.Id)
}

func (u *Index) SetId(v string) {
	u.Id = v
}

func (u *Index) GetTenantId() string {
	return u.TenantId
}

func (u *Index) SetTenantId(v string) {
	u.TenantId = v
}

func (u *Index) GetId() string {
	return string(u.Id)
}

func (u *Record) SetId(v string) {
	u.Id = v
}
