package element

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/runtime"
	"github.com/liuxd6825/dapr-go-ddd-sdk/os/fs/fsopts"
	"github.com/liuxd6825/jsonschema/v6"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type FuncConfig struct {
	FuncName    string                // 方法名称
	Code        string                // 代码内容
	CodeType    string                // 代码类型
	SrcFileName string                // 源文件名称
	UsePool     bool                  // 是否使用脚本缓存
	ParamsUrl   string                // 方法参数定义
	ParamsType  string                // 参数类型
	TransType   runtime.TransformType //转换类型
	Tags        []*FuncTag            // 其它扩展标签
	Selection   *goquery.Selection    // HTML选择器
	TimeoutAttr string                // 请求超时时间
	TxAttr      string                // 数据库事务HTML属性
	txDbKeys    []string              // 要开启的数据库事务
	timeout     time.Duration         // 超时时间
}

type Func interface {
	Logger() logrus.FieldLogger
	BuildCode() error
	Run(ctx context.Context, opts ...RunOptions) (res any, err error)
	Config() *FuncConfig
	AsJsFunc() any
	//GetParamsType(urlPars map[string]any, fsOpt *fsopts.Options) (paramsTypeFile string, paramsType common.ParamsType)
	GetParamsSchema(urlPars map[string]any, fsOpt *fsopts.Options) *jsonschema.Schema
}

type FuncManager interface {
	Add(fun Func) error
	Run(ctx context.Context, funcName string, runValues *ApiRunValues, checkHave bool, opts ...RunOptions) (res any, err error)
	Get(funcName string) (fun Func, ok bool)
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
		TxAttr:      scriptEl.AttrOr("tx", ""),
		TimeoutAttr: scriptEl.AttrOr("timeout", ""),
	}
	cfg.txDbKeys = cfg.getTxDbKeys()
	cfg.timeout = cfg.getTimeout()
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
	if f.TxAttr == "" {
		return []string{}
	}
	txDbKeys := strings.Split(f.TxAttr, ",")
	return txDbKeys
}

func (f *FuncConfig) Timeout() time.Duration {
	return f.timeout
}

func (f *FuncConfig) getTimeout() time.Duration {
	if f.TimeoutAttr == "" {
		return time.Duration(0)
	}
	timeout, err := time.ParseDuration(f.TimeoutAttr)
	if err != nil {
		panic(errors.ErrorOf("无效的时间格式:%s", err.Error()))
	}
	f.timeout = timeout
	return timeout
}
