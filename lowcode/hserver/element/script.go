package element

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/runtime"
	"github.com/sirupsen/logrus"
)

type ScriptConfig struct {
	FuncName    string
	Code        string
	CodeType    string
	SrcFileName string
	UsePool     bool
	Alias       map[string]string
	TransType   runtime.TransformType //转换类型
}

type Script interface {
	Logger() logrus.FieldLogger
	BuildCode() error
	Run(opts ...RunOptions) (res any, err error)
	Config() *ScriptConfig
}
