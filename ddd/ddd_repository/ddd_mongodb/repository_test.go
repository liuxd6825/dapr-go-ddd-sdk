package ddd_mongodb

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

func Test_Search(t *testing.T) {
	ctx := context.Background()
	mongodb, coll := newCollection("test_users")
	repository := newRepository(mongodb, coll)
	objId := NewObjectID()
	id := string(objId)
	user := &User{
		Id:        id,
		TenantId:  "001",
		UserName:  "UserName",
		UserCode:  "UserCode",
		Address:   "address",
		Email:     "lxd@163.com",
		Telephone: "17767788888",
	}

	_ = repository.Insert(ctx, user).OnSuccess(func(data *User) error {
		println(data)
		return nil
	}).OnError(func(err error) error {
		assert.Error(t, err)
		return err
	})

	search := ddd_repository.NewFindPagingQuery()
	search.SetTenantId("001")
	search.SetFilter(fmt.Sprintf("id=='%s'", id))

	_ = repository.FindPaging(ctx, search).OnSuccess(func(data []*User) error {
		println(data)
		return nil
	}).OnNotFond(func() error {
		err := errors.NewNotFondError()
		assert.Error(t, err)
		return err
	}).OnError(func(err error) error {
		assert.Error(t, err)
		return err
	}).GetError()

}

func TestMongoSession_UseTransaction(t *testing.T) {
	mongodb, coll := newCollection("test_users")
	repository := newRepository(mongodb, coll)
	err := ddd_repository.StartSession(context.Background(), NewSession(true, mongodb), func(ctx context.Context) error {
		for i := 0; i < 5; i++ {
			user := &User{
				Id:        NewObjectID().String(),
				TenantId:  "001",
				UserName:  "userName" + fmt.Sprint(i),
				UserCode:  "UserCode",
				Address:   "address",
				Email:     "lxd@163.com",
				Telephone: "17767788888",
			}

			err := repository.Insert(ctx, user).OnSuccess(func(data *User) error {
				println(data)
				return nil
			}).GetError()
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}
func TestRSqlIsNull(t *testing.T) {
	mongodb, coll := newCollection("record_ie")
	ctx := context.Background()
	repos := NewRepository[*RecordIe](func() *RecordIe { return &RecordIe{} }, mongodb, coll)
	qry := ddd_repository.NewFindPagingQuery()
	qry.SetFilter("errors=!null=0")
	qry.SetTenantId("test")
	qry.SetIsTotalRows(true)
	res := repos.FindPaging(ctx, qry)
	totalRow := *res.TotalRows
	t.Logf("totalRow:%v", totalRow)
	assert.Greater(t, totalRow, int64(0))
}
func newRepository(mongodb *MongoDB, coll *mongo.Collection) *Repository[*User] {
	return NewRepository[*User](func() *User { return &User{} }, mongodb, coll)
}

func newCollection(name string) (*MongoDB, *mongo.Collection) {
	config := &Config{
		Host:         "122.143.11.104:27018,122.143.11.104:27019,122.143.11.104:27020",
		DatabaseName: "duxm_fund_record_service_cmd_db",
		UserName:     "fund_record",
		Password:     "123456",
		Options:      "replicaSet=mongors",
	}
	mongodb, err := NewMongoDB(config)
	if err != nil {
		panic(err)
	}
	_ = mongodb.CreateCollection(name)
	coll := mongodb.GetCollection(name)
	return mongodb, coll
}

type User struct {
	Id        string `json:"id" validate:"gt=0" bson:"_id"`
	TenantId  string `json:"tenantId" validate:"gt=0" bson:"tenant_id"`
	UserCode  string `json:"userCode" validate:"gt=0" bson:"user_code"`
	UserName  string `json:"userName" validate:"gt=0" bson:"user_name"`
	Email     string `json:"email" validate:"gt=0" bson:"email"`
	Telephone string `json:"telephone" validate:"gt=0" bson:"telephone"`
	Address   string `json:"address" validate:"gt=0" bson:"address"`
}

func (u *User) GetTenantId() string {
	return u.TenantId
}

func (u *User) GetId() string {
	return string(u.Id)
}

func (u *User) SetTenantId(v string) {
	u.TenantId = v
}

func (u *User) SetId(v string) {
	u.Id = v
}

type RecordIe struct {
	Id     string `json:"id,omitempty"  bson:"_id"   validate:"required"  description:"id"` // 行Id
	RowNum int64  `json:"rowNum,omitempty" bson:"rowNum" description:"租户标识"`

	TenantId string `json:"tenantId,omitempty" bson:"tenant_id" description:"租户标识"`
	CaseId   string `json:"caseId,omitempty" bson:"case_id" description:"案件id"`
	TaskId   string `json:"taskId,omitempty" bson:"task_id"  description:"任务id"`
	DocId    string `json:"docId,omitempty" bson:"doc_id"  description:"文档id"`
	FileId   string `json:"fileId,omitempty" bson:"file_id"  description:"文件id"`

	Iden     string   `json:"iden,omitempty"  bson:"iden"   validate:"-" description:"我方标识"`            // 标识
	Name     string   `json:"name,omitempty"   bson:"name"  validate:"-" description:"我方名称"`            // 名称
	Acct     string   `json:"acct,omitempty"   bson:"acct"  validate:"-" description:"我方账号"`            // 账号
	AcctType string   `json:"acctType,omitempty"   bson:"acct_type"  validate:"-" description:"我方账号类型"` // 账号类型
	Category string   `json:"category,omitempty"   bson:"category"  validate:"-" description:"我方类别"`    // 类别Id 公司或个人
	BankName string   `json:"bankName,omitempty"  bson:"bank_name"   validate:"-" description:"我方开户银行"` // 开户银行
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

	Errors   map[string][]string `json:"errors,omitempty" bson:"errors,omitempty"`
	Cells    map[string]Cells    `json:"cells,omitempty" bson:"cells,omitempty"`
	HasError bool                `json:"hasError" bson:"hasError,omitempty"`
}

type Cells []Cell

type Cell struct {
	Key    string   `json:"key" bson:"key,omitempty"`
	Value  string   `json:"value" bson:"value,omitempty"`
	Col    int64    `json:"col" bson:"col,omitempty"`
	Row    int64    `json:"row" bson:"row,omitempty"`
	Errors []string `json:"errors" bson:"errors,omitempty"`
}

func (u *RecordIe) GetTenantId() string {
	return u.TenantId
}

func (u *RecordIe) GetId() string {
	return string(u.Id)
}

func (u *RecordIe) SetTenantId(v string) {
	u.TenantId = v
}

func (u *RecordIe) SetId(v string) {
	u.Id = v
}
