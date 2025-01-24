package element

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/common"
)

func ParseScriptConfig(parentEl *goquery.Selection, selector string, funcName string, srcFileName string) (*ScriptConfig, error) {
	var err error
	var cfg *ScriptConfig
	if parentEl == nil {
		return cfg, nil
	}
	scripts := parentEl.Find(selector).First()
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
