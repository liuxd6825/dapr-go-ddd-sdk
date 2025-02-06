package element

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/runtime"
	"github.com/sirupsen/logrus"
)

type FuncConfig struct {
	FuncName    string
	Code        string
	CodeType    string
	SrcFileName string
	UsePool     bool
	Params      string                // 方法参数定义
	TransType   runtime.TransformType //转换类型
	Tags        []*FuncTag
	Selection   *goquery.Selection
}

type Func interface {
	Logger() logrus.FieldLogger
	BuildCode() error
	Run(opts ...RunOptions) (res any, err error)
	Config() *FuncConfig
}

type FuncManager interface {
	Add(fun Func) error
	Run(funcName string, runValues *ApiRunValues, checkHave bool, opts ...RunOptions) (res any, err error)
	Items() map[string]Func
}

func NewFuncConfig(scriptEl *goquery.Selection, srcFileName string, tags []*FuncTag) (cfg *FuncConfig, isFound bool, err error) {
	if scriptEl == nil {
		return cfg, false, nil
	}
	node := scriptEl.Get(0)
	if node.Data != "script" {
		return cfg, false, errors.New("script node is not script")
	}

	cfg = &FuncConfig{
		Selection:   scriptEl,
		Tags:        tags,
		SrcFileName: srcFileName,
		Params:      scriptEl.AttrOr("params", ""),
		FuncName:    scriptEl.AttrOr("name", ""),
		CodeType:    scriptEl.AttrOr("type", ""),
		Code:        scriptEl.Text(),
	}

	return cfg, true, err
}

func (f *FuncConfig) ContainTag(tagType string) bool {
	for _, tag := range f.Tags {
		if tag.Type == tagType {
			return true
		}
	}
	return false
}
