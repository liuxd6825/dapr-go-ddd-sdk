package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/enum"
)

type Tenant struct {
	Base    `bson:",inline"`
	Name    string      `json:"name" gorm:"name" bson:"name" `
	Link    string      `json:"link" bson:"link" gorm:"link"`
	Address string      `json:"address" gorm:"address" bson:"address"`
	Phone   string      `json:"phone" gorm:"phone" bson:"phone"`
	Email   string      `json:"email" gorm:"email" bson:"email"`
	Status  enum.Status `json:"status" gorm:"status" bson:"status"` // 状态
}

func NewTenant() (*Tenant, error) {
	return &Tenant{}, nil
}
