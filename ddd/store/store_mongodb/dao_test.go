package store_mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/randomutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
	"time"
)

const DB_NAME = "test"
const TENANT_ID = "test"

func newHumanDao() store.IStore[*Human] {
	humanDbSch := dbschema.NewDBSchemaWithStruct("human", &Human{}, "human")
	humanDao := NewDao[*Human](humanDbSch, func(ctx context.Context) (IMongoDB, *mongo.Collection) {
		return mongodb.NenMongoDBWithClient(DB_NAME, client), getCollection(DB_NAME, "human")
	})
	return humanDao
}

func newRecordDao() store.IStore[*Record] {
	recordDbSch := dbschema.NewDBSchemaWithStruct("record", &Record{}, "record")
	recordDao := NewDao[*Record](recordDbSch, func(ctx context.Context) (IMongoDB, *mongo.Collection) {
		return mongodb.NenMongoDBWithClient(DB_NAME, client), getCollection(DB_NAME, "record")
	})
	return recordDao
}

func TestMapper_InsertUpdate(t *testing.T) {
	ctx := xtest.NewContext()
	humanDao := newHumanDao()
	id := randomutils.NewId()
	human := &Human{
		Id:       id,
		Name:     "TEST-A",
		TenantId: TENANT_ID,
	}
	humanDao.Insert(ctx, human)

	human.Name = "TEST-B"
	opts := store.NewOptions()
	opts.SetUpdateFields([]string{"name"})
	humanDao.Update(ctx, human, opts)

	result := humanDao.FindById(ctx, TENANT_ID, id)
	if result.GetError() != nil {
		t.Fatal(result.GetError().Error())
	}

	assert.Equal(t, human.Name, "TEST-B")

}

func TestMapper_FindById(t *testing.T) {
	ctx := xtest.NewContext()
	humanDao := newHumanDao()
	result := humanDao.FindById(ctx, TENANT_ID, "SKCZO381BFONRJRUTD8T7M30XYISP9")
	assert.NotNil(t, result.Data)
}

func TestMapper_Search(t *testing.T) {
	ctx := xtest.NewContext()

	humanDao := newHumanDao()
	recordDao := newRecordDao()

	humanName := "张三"
	t.Run("Inserts", func(t *testing.T) {
		humanMax := 1
		var humanList []*Human
		var recordList []*Record
		for i := 0; i < humanMax; i++ {
			human := &Human{
				Id:       idutils.NewId(),
				Name:     humanName,
				TenantId: TENANT_ID,
			}
			humanList = append(humanList, human)
		}

		for i := 0; i < humanMax; i++ {
			human := humanList[i]
			for j := 0; j < 10; j++ {
				record := &Record{
					Id:       idutils.NewId(),
					Name:     human.Name,
					CaseId:   human.TenantId,
					TenantId: TENANT_ID,
					Acct:     randomutils.StringNumber(10),
					Iden:     randomutils.StringNumber(10), // 标识
					OppIden:  randomutils.StringNumber(10),
					OppName:  randomutils.StringNumber(10),
				}
				recordList = append(recordList, record)
			}

		}

		err1 := humanDao.InsertMany(ctx, TENANT_ID, humanList).Error
		assert.Nil(t, err1)

		err2 := recordDao.InsertMany(ctx, TENANT_ID, recordList).Error
		assert.Nil(t, err2)
	})

	t.Run("FindPagingQuery_NoGroup", func(t *testing.T) {
		qry := store.NewFindPagingQuery()
		filter := fmt.Sprintf("name=='%s' and tenantId=='%s'", humanName, TENANT_ID)
		qry.SetFilter(filter)
		qry.SetTenantId(TENANT_ID)
		qry.SetPageSize(20)

		res := recordDao.FindPaging(ctx, qry)
		if res.GetError() == nil {
			logObject(t, "Data: ", res.GetData())
			logObject(t, "Sum: ", res.GetSumData())
		} else {
			t.Error(res.GetError())
		}
	})

	t.Run("FindPagingQuery_Group", func(t *testing.T) {
		qry := store.NewFindPagingQuery()
		filter := fmt.Sprintf("name=='%s'", humanName)
		qry.SetFilter(filter)
		qry.SetTenantId(TENANT_ID)
		qry.SetPageSize(20)
		qry.SetGroupCols(store.NewGroupCols("").Add("name", types.DataTypeString).GetCols())
		qry.SetGroupKeys([]any{humanName})

		res := recordDao.FindPaging(ctx, qry)
		if res.GetError() == nil {
			logObject(t, "Data: ", res.GetData())
			logObject(t, "Sum: ", res.GetSumData())
		} else {
			t.Error(res.GetError())
		}
	})

	t.Run("Filter_Sub", func(t *testing.T) {
		rSql := fmt.Sprintf("name==sub(table:human, field:name, rsql:name~='%s')", humanName)
		res := recordDao.FindByRSQL(ctx, TENANT_ID, rSql)
		if res.Error == nil {
			logObject(t, "Data: ", res.Data)
		} else {
			t.Error(res.Error)
		}
	})

	t.Run("Filter_Like", func(t *testing.T) {
		rSql := fmt.Sprintf("name~='%s'", humanName)
		res := recordDao.FindByRSQL(ctx, TENANT_ID, rSql)
		if res.Error == nil {
			logObject(t, "Data: ", res.Data)
		} else {
			t.Error(res.Error)
		}
	})
}

func TestDao_CreateIndexes(t *testing.T) {
	/*
		ctx := context.Background()
		coll := getCollection(DB_NAME, "test_create_index")
		mapper := NewDao[*Index](nil, func(ctx context.Context) (IMongoDB, *mongo.Collection) {
			return mongodb.NenMongoDBWithClient(DB_NAME, client), coll
		})


		err := mapper.CreateIndexes(ctx)

		if err != nil {
			t.Error(err)
		}*/
}

func logObject(t *testing.T, label string, obj any) {
	jsonText, err := json.Marshal(obj)
	if err != nil {
		t.Error(label, err)
	}
	t.Log(label, string(jsonText))
}

type Index struct {
	Id        string `bson:"id" `
	TenantId  string
	Name      string `bson:"name" index:"" `
	Asc       int64  `bson:"asc" index:" asc"`
	Desc      int64  `bson:"desc" index:" desc "`
	Unique    string `index:"unique"`
	AscUnique string `bson:"asc_unique" index:"asc, unique "`
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

func (u *Index) SetId(v string) {
	u.Id = v
}

type Human struct {
	Id       string `bson:"id" `
	Name     string `bson:"name" index:"" `
	TenantId string `bson:"tenant_id" index:"" `
}

func (u *Human) GetTenantId() string {
	return u.TenantId
}

func (u *Human) SetTenantId(v string) {
	u.TenantId = v
}

func (u *Human) GetId() string {
	return string(u.Id)
}

func (u *Human) SetId(v string) {
	u.Id = v
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

func (u *Record) SetId(v string) {
	u.Id = v
}

var client *mongo.Client

func init() {
	// 设置MongoDB连接URL
	opts := options.Client().ApplyURI("mongodb://127.0.0.1:27017,127.0.0.1:27018,127.0.0.1:27019/?retryWrites=false&replicaSet=rs0&readPreference=primary&serverSelectionTimeoutMS=5000&connectTimeoutMS=10000")
	opts.Auth = &options.Credential{
		AuthMechanism: "SCRAM-SHA-256",
		AuthSource:    "admin",
		Username:      "super_admin",
		Password:      "123456",
	}
	// 连接到MongoDB
	clientVal, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}
	client = clientVal
	// 确保连接成功
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getCollection(dbName string, collName string) *mongo.Collection {
	// 选择数据库和集合
	collection := client.Database(dbName).Collection(collName)
	return collection
}
