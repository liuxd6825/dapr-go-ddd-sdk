package model

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"

type Company struct {
	xbase.BaseModel `bson:",inline"`
}
