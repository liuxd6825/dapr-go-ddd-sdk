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

func Test_NewObject(t *testing.T) {
	v, err := NewObject[map[string]any]()
	assert.NoError(t, err)
	t.Log(v)
}

type ActionLogProxy = ActionLog

func Test_GetTypeDetails(t *testing.T) {
	pkgPath, typeName, fields := GetTypeDetails(&ActionLog{})
	t.Logf("pkgPath: %s, typeName: %s, fields: %v", pkgPath, typeName, fields)

	pkgPath, typeName, fields = GetTypeDetails(&ActionLogProxy{})
	t.Logf("pkgPath: %s, typeName: %s, fields: %v", pkgPath, typeName, fields)
}
