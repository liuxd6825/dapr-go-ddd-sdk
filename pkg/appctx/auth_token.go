package appctx

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/jsonutils"
)

type AuthToken interface {
	GetSub() string
	GetExp() int
	GetUser() AuthUser
	GetClientId() string
	GetToken() string
	Copy(source AuthToken)
}

type AuthTokenEntity struct {
	Sub      string          `json:"sub"`
	Exp      int             `json:"exp"`
	User     *AuthUserEntity `json:"user"`
	ClientId string          `json:"client_id"`
	Token    string          `json:"token"`
}

func getAuthToken(tokenStr string) (AuthToken, error) {
	list := strings.Split(tokenStr, ".")
	if len(list) != 3 {
		return nil, errors.New("token格式不正确")
	}
	tk := &AuthTokenEntity{}
	bs, err := base64.RawURLEncoding.DecodeString(list[1])
	if err != nil {
		return nil, err
	}
	err = jsonutils.Unmarshal(bs, tk)
	tk.Token = tokenStr
	return tk, err
}

///////////////////////
//     authToken
///////////////////////

func NewAuthToken() *AuthTokenEntity {
	return &AuthTokenEntity{
		User: NewAuthUser(),
	}
}

func (u *AuthTokenEntity) Copy(source AuthToken) {
	u.Exp = source.GetExp()
	u.User = source.GetUser().(*AuthUserEntity)
	u.Sub = source.GetSub()
	u.ClientId = source.GetClientId()
}

func (u *AuthTokenEntity) GetSub() string {
	return u.Sub
}

func (u *AuthTokenEntity) GetExp() int {
	return u.Exp
}

func (u *AuthTokenEntity) GetUser() AuthUser {
	return u.User
}

func (u *AuthTokenEntity) GetClientId() string {
	return u.ClientId
}

func (u *AuthTokenEntity) GetToken() string {
	return u.Token
}
