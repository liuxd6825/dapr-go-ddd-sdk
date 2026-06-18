package model

type IOType int

const (
	IOType_Out IOType = 0
	IOType_In  IOType = 1
)

type CashType int

const (
	CashType_None CashType = 0 // 非现金
	CashType_Cash CashType = 1 // 现金
)

//type CashType bool
//
//const (
//	CashType_None CashType = false // 非现金
//	CashType_Cash CashType = true  // 现金
//)

type AccountType string

const (
	AccountType_Personal AccountType = "个人"
	AccountType_Company  AccountType = "公司"
)
