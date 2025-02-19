package element

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/runtime"
	"github.com/sirupsen/logrus"
	"strings"
)

type FuncConfig struct {
	FuncName    string
	Code        string
	CodeType    string
	SrcFileName string
	UsePool     bool
	ParamsUrl   string // 方法参数定义
	ParamsType  string
	TransType   runtime.TransformType //转换类型
	Tags        []*FuncTag
	Selection   *goquery.Selection
	Tx          string   // 数据库事务HTML属性
	txDbKeys    []string // 要开启的数据库事务
}

type Func interface {
	Logger() logrus.FieldLogger
	BuildCode() error
	Run(opts ...RunOptions) (res any, err error)
	Config() *FuncConfig
	GetParamsType(urlPars map[string]any, fsOpt *fsopts.Options) (paramsTypeFile string, paramsType common.ParamsType)
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
		ParamsUrl:   scriptEl.AttrOr("params-url", ""),
		ParamsType:  scriptEl.AttrOr("params-type", ""),
		FuncName:    scriptEl.AttrOr("name", ""),
		CodeType:    scriptEl.AttrOr("type", ""),
		Code:        scriptEl.Text(),
		Tx:          scriptEl.AttrOr("tx", ""),
	}
	cfg.txDbKeys = cfg.getTxDbKeys()

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

func (f *FuncConfig) TxDbKeys() []string {
	return f.txDbKeys
}

func (f *FuncConfig) getTxDbKeys() []string {
	if f.Tx == "" {
		return []string{}
	}
	txDbKeys := strings.Split(f.Tx, ",")
	return txDbKeys
}
