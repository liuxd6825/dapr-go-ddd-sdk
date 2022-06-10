package mapper

import (
	"encoding/json"
	"testing"
)

type UserCreateRequest struct {
	CommandId   string                `json:"commandId"`
	IsValidOnly bool                  `json:"isValidOnly"`
	Data        UserCreateRequestData `json:"data"`
}

type UserCreateRequestData struct {
	UserRequestData
}

type UserRequestData struct {
	Id        string `json:"id" validate:"gt=0" minLength:"16" maxLength:"16" example:"random string"`
	TenantId  string `json:"tenantId" validate:"gt=0" minLength:"16" maxLength:"16" example:"random string"`
	UserCode  string `json:"userCode" validate:"gt=0"  minLength:"16" maxLength:"16" example:"random string"`
	UserName  string `json:"userName" validate:"gt=0"  minLength:"16" maxLength:"16" example:"random string"`
	Email     string `json:"email" validate:"gt=0"    example:"xxx@163.com"`
	Telephone string `json:"telephone" validate:"gt=0"  length:"11" example:"18867766829"`
	Address   string `json:"address" validate:"gt=0"`
}

type UserCreateAppCommand struct {
	UserCreateCommand
}

type UserCreateCommand struct {
	CommandId   string     `json:"commandId"  validate:"gt=0"`
	IsValidOnly bool       `json:"isValidOnly"`
	Data        UserFields `json:"data"`
}

type UserFields struct {
	Id        string `json:"id" validate:"gt=0"`
	TenantId  string `json:"tenantId" validate:"gt=0"`
	UserCode  string `json:"userCode" validate:"gt=0"`
	UserName  string `json:"userName" validate:"gt=0"`
	Email     string `json:"email" validate:"gt=0"`
	Telephone string `json:"telephone" validate:"gt=0"`
	Address   string `json:"address" validate:"gt=0"`
}

func TestAutoMapper(t *testing.T) {
	request := UserCreateRequest{
		CommandId:   "000",
		IsValidOnly: true,
		Data: UserCreateRequestData{
			UserRequestData: UserRequestData{
				Id:        "1111",
				Telephone: "1112222",
				Address:   "address",
			},
		},
	}
	var command UserCreateAppCommand
	if err := Mapper(&request, &command); err != nil {
		t.Error(err)
		return
	}
	jsonByte, _ := json.Marshal(command)
	print(string(jsonByte))
	if command.Data.Id == "" {
		t.Error("command.data.id is error")
	}

}
