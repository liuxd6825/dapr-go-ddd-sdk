package mapper

import (
	"encoding/json"
	"errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"testing"
	"time"
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

type DateRequest struct {
	Date *types.JSONDate
}

type DateCommand struct {
	Date time.Time
}

func TestDateMapper(t *testing.T) {
	dateValue := types.JSONDate(time.Now())
	// dateValue := types.DateString("2019-10-10")
	// dateValue := time.Now()
	req := DateRequest{
		Date: &dateValue,
	}
	cmd := DateCommand{}
	if err := Mapper(&req, &cmd); err != nil {
		t.Error(err)
	}
	if cmd.Date.IsZero() {
		t.Error(errors.New(" date mapper error"))
	}
	println(cmd.Date.String())
}

func TestMapperMask(t *testing.T) {
	from := UserFields{
		Id:       "0001",
		TenantId: "tenantId",
		UserName: "userName",
	}
	to := UserFields{UserName: "userName___"}
	mask := []string{
		"Id",
		"TenantId",
	}
	if err := MaskMapper(&from, &to, mask); err != nil {
		t.Error(err)
	}
	if to.UserName != "" {
		t.Error(errors.New("to.UserName is not null"))
	}
	if to.Id == "" {
		t.Error(errors.New("to.Id is null"))
	}
}
