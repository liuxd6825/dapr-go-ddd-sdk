package model

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"

type Account struct {
	xbase.BaseModel `bson:",inline"`
	Account         string `json:"account" bson:"account" title:"账号"`
	BankName        string `json:"bankName" bson:"bank_name" title:"银号"`
	OwnerName       string ` json:"ownerName" bson:"owner_name" title:"持卡人"`
}
