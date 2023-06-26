package readexcel

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testReviewFile(t *testing.T, file string, sheetName string, maxRows int) {
	review, err := ReviewFile(file, sheetName, maxRows)
	h := ""
	for i := 0; i < len(review.Columns); i++ {
		h += review.Columns[i] + " , "
	}
	t.Log(h)

	for _, item := range review.Items {
		s := ""
		for i := 0; i < len(review.Columns); i++ {
			key := review.Columns[i]
			s += item[key] + ","
		}
		t.Log(s)
	}
	if err != nil {
		t.Error(err)
	}
}

func TestOpen(t *testing.T) {
	testReviewFile(t, "./test_files/100w.xlsx", "Sheet0", 100000)
}

func TestReviewFile(t *testing.T) {
	testReviewFile(t, "./test_files/record.xlsx", "Sheet2", 100)
}

func TestReviewFile_100W(t *testing.T) {
	testReviewFile(t, "./test_files/100w.xlsx", "Sheet0", 100000)
}

func TestReadFile_Record(t *testing.T) {
	tmp, err := newRecordTemp()
	if err != nil {
		t.Error(err)
		return
	}
	if table, err := ReadFile("./test_files/record.xlsx", "Sheet2", tmp); err != nil {
		t.Error(err)
	} else {
		for _, row := range table.Rows {
			if b, err := json.Marshal(row.Values); err != nil {
				t.Error(err)
				t.Log("rowNum:", row.RowNum)
			} else {
				t.Log(string(b))
			}
			fmt.Println()
		}
	}
}

func TestReadFile_70W(t *testing.T) {
	tmp, err := new70wTemp()
	if err != nil {
		t.Error(err)
		return
	}
	if table, err := ReadFile("./test_files/70w.xlsx", "Sheet1", tmp); err != nil {
		t.Error(err)
	} else {
		t.Log(len(table.Rows))
	}
}

func TestGetCellLabel(t *testing.T) {
	type testItem struct {
		label string
		value int
	}
	list := []testItem{
		{
			label: "A",
			value: 1,
		}, {
			label: "Z",
			value: 26,
		}, {
			label: "BZ",
			value: 78,
		}, {
			label: "CA",
			value: 79,
		}, {
			label: "DT",
			value: 124,
		}, {
			label: "CZ",
			value: 104,
		}, {
			label: "DA",
			value: 105,
		},
	}
	for _, item := range list {
		_ = assert.Equal(t, item.label, GetCellLabel(item.value))
	}
}

func new70wTemp() (*Template, error) {
	fields := []*Field{
		NewField("f1", "商户名称", String).SetMapKeys("商户名称"),
		NewField("f2", "商户编号", String).SetMapKeys("商户编号"),
		NewField("f3", "交易账号", String).SetMapKeys("交易账号"),
		NewField("f4", "交易户名", String).SetMapKeys("交易户名"),
		NewField("f5", "交易账户开户银行", String).SetMapKeys("交易账户开户银行"),
		NewField("f6", "交易日期", String).SetMapKeys("交易日期"),
		NewField("f7", "交易类型", String).SetMapKeys("交易类型"),
		NewField("f8", "交易金额", String).SetMapKeys("交易金额"),
	}

	heads := []*MapHead{
		{
			RowNum: 1,
			Columns: []*MapColumn{
				{Key: "商户名称", Label: "A"}, // 交易时间
				{Key: "商户编号", Label: "B"}, // 我方账号
				{Key: "交易账号", Label: "C"},
				{Key: "交易户名", Label: "D"},
				{Key: "交易账户开户银行", Label: "E"},
				{Key: "交易日期", Label: "F"},
				{Key: "交易类型", Label: "G"},
				{Key: "交易金额", Label: "H"},
			},
		},
	}
	return NewTemplate(fields, heads, nil)
}

func newRecordTemp() (*Template, error) {
	fields := []*Field{
		NewField("Iden", "标识", String),
		NewField("Acct", "账号", String).SetMapKeys("Acct"),
		NewField("Name", "名称", String).SetMapKeys("Name"),
		NewField("AcctType", "账号类型", String),
		NewField("BankName", "开户行", String).SetMapKeys("BankName"),
		NewField("Balance", "余额", String).SetMapKeys("Balance"),
		NewField("Category", "类别", String),

		NewField("OppAcct", "对方账号", String).SetMapKeys("OppAcct"),
		NewField("OppAcctType", "对方账号类型", String),
		NewField("OppCategory", "对方类别", String),
		NewField("OppIden", "对方标识", String),
		NewField("OppBankName", "对方开户行", String).SetMapKeys("OppBankName"),
		NewField("OppName", "对方名称", String).SetMapKeys("OppName"),

		NewField("Amount", "交易金额", Money).SetMapKeys("Income", "Payout").SetScript(`Income!='0.00'?Income:Payout`),
		NewField("Date", "交易时间", DateTime).SetMapKeys("Date"),
		NewField("Notes", "交易备注", String).SetMapKeys("Summary", "Purpose", "Remarks"),
		NewField("Place", "交易地点", String),
		NewField("Ccy", "交易币种", String), //.SetDefault(DefaultValue("RMB")),

		NewField("Serial", "流水号", String),
		NewField("Income", "收入金额", String).SetMapKeys("Income"),
		NewField("Payout", "支出金额", String).SetMapKeys("Payout"),
	}

	heads := []*MapHead{
		{
			RowNum: 4,
			Columns: []*MapColumn{
				{Key: "Date", Label: "A"}, // 交易时间
				{Key: "Acct", Label: "B"}, // 我方账号
				{Key: "Name", Label: "C"},
				{Key: "BankName", Label: "D"},
				{Key: "OppAcct", Label: "E"},
				{Key: "OppName", Label: "F"},
				{Key: "OppBankName", Label: "E"},
				{Key: "Payout", Label: "H"},
				{Key: "Income", Label: "I"},
				{Key: "Balance", Label: "J"},
				{Key: "Summary", Label: "K"},
				{Key: "Purpose", Label: "L"},
				{Key: "Remarks", Label: "M"},
			},
		},
	}

	consts := []*Const{
		{Key: "MyName", Type: ConstTypeValue, Value: PString("北京家家财富投资有限公司")},
		{Key: "MyAccount", Type: ConstTypeCell, Point: &Point{Row: 1, Col: 1}},
	}

	return NewTemplate(fields, heads, consts)
}
