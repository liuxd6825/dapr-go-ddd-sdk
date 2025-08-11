package model

type IOType int

const (
	IOType_Out IOType = 0
	IOType_In  IOType = 1
)

type CashType bool

const (
	CashType_None CashType = false // 非现金
	CashType_Cash CashType = true  // 现金
)

type AcctType string

const (
	AcctType_Personal AcctType = "个人"
	AcctType_Company  AcctType = "公司"
)
