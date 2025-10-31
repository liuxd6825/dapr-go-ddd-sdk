package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/enum"
	"time"
)

type Tenant struct {
	Base       `bson:",inline"`
	Name       string            `json:"name" gorm:"name" bson:"name"`
	Code       string            `json:"code" gorm:"code" bson:"code"`
	Address    string            `json:"address" gorm:"address" bson:"address"`
	Phone      string            `json:"phone" gorm:"phone" bson:"phone"`
	Email      string            `json:"email" gorm:"email" bson:"email"`
	EffectDate *time.Time        `json:"effectDate" gorm:"effect_date" bson:"effect_date"` //生效日期
	ExpireDate *time.Time        `json:"expireDate" gorm:"expire_date" bson:"expire_date"` //过期日期
	LogoutDate *time.Time        `json:"logoutDate" gorm:"logout_date" bson:"logout_date"` //注销日期
	Status     enum.TenantStatus `json:"status" gorm:"status" bson:"status"`               // 状态
}

func NewTenant() (*Tenant, error) {
	return &Tenant{}, nil
}
