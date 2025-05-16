package appctx

type AuthUser interface {
	GetId() string
	GetName() string
	GetPhone() string

	GetAccount() string
	GetRegDate() string
	GetWork() string
	GetStatus() string
	GetUserType() string

	GetTenantId() string
	GetTenantName() string
}

type AuthUserEntity struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Account    string `json:"account"`
	RegDate    string `json:"regDate"`
	Work       string `json:"work"`
	Status     string `json:"status"`
	UserType   string `json:"userType"`
	TenantId   string `json:"tenantId"`
	TenantName string `json:"tenantName"`
}

func NewAuthUser() *AuthUserEntity {
	return &AuthUserEntity{}
}

func (a *AuthUserEntity) GetId() string {
	return a.Id
}

func (a *AuthUserEntity) GetName() string {
	return a.Name
}

func (a *AuthUserEntity) GetPhone() string {
	return a.Phone
}

func (a *AuthUserEntity) GetAccount() string {
	return a.Account
}

func (a *AuthUserEntity) GetRegDate() string {
	return a.RegDate
}

func (a *AuthUserEntity) GetWork() string {
	return a.Work
}

func (a *AuthUserEntity) GetStatus() string {
	return a.Status
}

func (a *AuthUserEntity) GetUserType() string {
	return a.UserType
}

func (a *AuthUserEntity) GetTenantId() string {
	return a.TenantId
}

func (a *AuthUserEntity) GetTenantName() string {
	return a.TenantName
}
