package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

// SuTaskAccount 分析的账户
type SuTaskAccount struct {
	xbase.BaseModel `bson:",inline"`
	TaskId          string            `json:"taskId" gorm:"task_id" bson:"task_id"`
	Account         string            `json:"account" gorm:"account"  bson:"account" title:"账户"`
	AccountType     model.AccountType `json:"accountType" gorm:"account_type"  bson:"account_type" title:"账户类型"`
	Name            string            `json:"name" gorm:"name"  bson:"name" title:"账户名称"`
	BankName        string            `json:"bankName" gorm:"bank_name" bson:"bank_name" title:"开户银行"`
}
