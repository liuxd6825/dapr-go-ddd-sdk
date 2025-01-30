package element

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/runtime"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
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

type ScriptManager interface {
	AddScript(config *ScriptConfig, logger logrus.FieldLogger, pkg *types.CMap[any]) error
	RunScript(funcName string, runValues *ApiRunValues, checkHave bool, opts ...RunOptions) (res any, err error)
}

func GetScriptConfig(parentEl *goquery.Selection, funcName string, srcFileName string) (cfg *ScriptConfig, isFound bool, err error) {
	if parentEl == nil {
		return cfg, false, nil
	}
	sel := parentEl.Find("script").First()
	leg := sel.Length()
	if leg == 1 {
		cfg = &ScriptConfig{
			FuncName:    funcName,
			SrcFileName: srcFileName,
		}
		cfg.Code = sel.Text()
		cfg.CodeType = sel.AttrOr("type", "")
		alias := sel.AttrOr("alias", "")
		cfg.Alias, err = common.ParseAlias(alias)
	} else if leg > 1 {
		return cfg, false, errors.ErrNotFound
	} else if leg == 0 {
		return nil, false, nil
	}
	return cfg, false, err
}

func NewScriptConfig1(parentEl *goquery.Selection, funcName string, srcFileName string) (*ScriptConfig, error) {
	var err error
	var cfg *ScriptConfig
	if parentEl == nil {
		return cfg, nil
	}
	scripts := parentEl.Find("script").First()
	if scripts.Length() > 0 {
		cfg = &ScriptConfig{
			FuncName:    funcName,
			SrcFileName: srcFileName,
		}
		cfg.Code = scripts.Text()
		cfg.CodeType = scripts.AttrOr("type", "")
		alias := scripts.AttrOr("alias", "")
		cfg.Alias, err = common.ParseAlias(alias)
	} else {
		return nil, errors.New("<script> must have a selector")
	}
	return cfg, err
}
