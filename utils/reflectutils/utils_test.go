package reflectutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type ActionLog struct {
	Id         string    `desc:"ID" ulog:"-"`
	ActionType string    `desc:"操作类型"`
	AppId      string    `desc:"应用ID"`
	AppName    string    `desc:"应用名称"`
	UserId     string    `desc:"用户ID"`
	UserName   string    `desc:"用户名称"`
	TenantId   string    `desc:"租户ID"`
	Time       time.Time `desc:"操作时间"`
	ModelName  string    `desc:"实体"`
	Message    string    `desc:"消息" `
	IntVal     int       ``
}

func TestParseDesc(t *testing.T) {
	desc, err := ParseDesc(&ActionLog{Id: "0001", Time: time.Now(), IntVal: 1000})
	assert.NoError(t, err)
	t.Logf("%s", desc)
}

func TestGetTypeName(t *testing.T) {

}
