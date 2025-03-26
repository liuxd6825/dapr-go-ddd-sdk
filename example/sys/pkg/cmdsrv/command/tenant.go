package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/example/sys/pkg/xinfra/enums"
	"time"
)

type TenantCommand struct {
	EffectDate    time.Time          `gorm:"column:effect_date" json:"effectDate"`
	LogoutDate    time.Time          `gorm:"column:logout_date" json:"logoutDate"`
	ExpireDate    time.Time          `gorm:"column:expire_date" json:"expireDate"`
	TenantStatus  enums.TenantStatus `gorm:"column:tenant_status" json:"tenantStatus"`
	TenantAddress string             `gorm:"column:tenant_address" json:"tenantAddress"`
	TenantEmail   string             `gorm:"column:tenant_email" json:"tenantEmail"`
	TenantName    string             `gorm:"column:tenant_name" json:"tenantName"`
	TenantPhone   string             `gorm:"column:tenant_phone" json:"tenantPhone"`
	TenantType    enums.TenantType   `gorm:"column:tenant_type" json:"tenantType"`
	TenantAccount string             `gorm:"column:tenant_account" json:"tenantAccount"`
	MenusID       []string           `gorm:"column:menus_id" json:"menusId"`
}
