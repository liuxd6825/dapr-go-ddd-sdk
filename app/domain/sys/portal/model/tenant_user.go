package model

type TenantUser struct {
	Base    `bson:",inline"`
	TenId   string `json:"tenId" bson:"ten_id" gorm:"ten_id"`
	UserId  string `json:"userId" bson:"user_id" gorm:"user_id"`
	IsAdmin bool   `json:"isAdmin" bson:"is_admin" gorm:"is_admin"`
}

func NewTenantUser() (*TenantUser, error) {
	return &TenantUser{}, nil
}
