package funcs

import (
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/runtime"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type Func struct {
	server         element.Server
	runCode        string
	runtime        *runtime.Pool
	config         *element.FuncConfig
	fs             afero.Fs
	transType      runtime.TransformType //转换类型
	isBuildCode    bool
	logger         logrus.FieldLogger
	paramsTypeFile string            // 参数文件
	paramsType     common.ParamsType // 参数类型定义
}

type RunOptions = func(vm *goja.Runtime) error

type InitVM interface {
	InitVM(vm *goja.Runtime) error
}

func NewFunc(server element.Server, config *element.FuncConfig, logger logrus.FieldLogger, reader fs.Reader, pkg *types.CMap[any]) (element.Func, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	fun := &Func{
		server:    server,
		runtime:   runtime.NewPool(config.UsePool, reader, pkg),
		config:    config,
		logger:    logger,
		transType: config.TransType,
	}

	err := fun.BuildCode()
	if err != nil {
		panic(err)
	}

	return fun, nil
}

func (s *Func) Logger() logrus.FieldLogger {
	return s.logger
}

func (s *Func) Config() *element.FuncConfig {
	return s.config
}

func (s *Func) BuildCode() error {
	if !s.isBuildCode {
		codeBytes, err := runtime.TransformCode(s.config.Code, s.config.SrcFileName, s.config.TransType)
		if err != nil {
			return err
		}
		s.isBuildCode = true
		s.runCode = string(codeBytes)
	}
	return nil
}

func (s *Func) Run(opts ...RunOptions) (res any, err error) {
	defer func() {
		err = utils.RecoverError(err, recover())
	}()

	if s != nil {
		opts = append(opts, func(vm *goja.Runtime) error {
			return nil
		})
	}

	data, err := s.runtime.Run("\n"+s.runCode, opts...)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("run error:%v in %s %s", err, s.config.FuncName, s.config.SrcFileName))
	}
	//r.logger.Printf("run %s return %v in %s", r.config.FuncName, data, r.config.SrcFileName)
	return data, err
}

// GetParamsType
//
//	@Description: 获取参数类型定义
//	@receiver r
//	@param ictx
//	@return string  参数文件名
//	@return common.ParamsType  参数配置类型
func (s *Func) GetParamsType(urlPars map[string]any, fsOpt *fsopts.Options) (paramsTypeFile string, paramsType common.ParamsType) {
	if s.paramsTypeFile != "" {
		return s.paramsTypeFile, s.paramsType
	}
	// 引用Schema文件
	if s.config.ParamsUrl != "" {
		paramsTypeFile, paramsType = s.getParamsTypeByUrl(s.config.ParamsUrl, urlPars, fsOpt)
	} else if s.config.ParamsType != "" {
		paramsTypeFile, paramsType = s.getParamsTypeByType(s.config.ParamsType, fsOpt)
	} else {
		paramsTypeFile, paramsType = s.getParamsTypeByUrl("./params/"+s.config.FuncName, urlPars, fsOpt)
		if paramsType == nil {
			paramsTypeFile, paramsType = s.getParamsTypeByType(s.config.FuncName, fsOpt)
		}
	}
	s.paramsTypeFile = paramsTypeFile
	s.paramsType = paramsType

	return s.paramsTypeFile, s.paramsType
}

func (s *Func) getParamsTypeByType(aParamsType string, fsOpt *fsopts.Options) (paramsTypeFile string, paramsType common.ParamsType) {
	fileName := fmt.Sprintf("/definition/params/%s.json", aParamsType)
	if pt := s.server.Definition().GetParamsType(aParamsType + ".json"); pt != nil {
		// 引用系统中的schema定义文件
		paramsTypeFile = fileName
		paramsType = pt
	}
	return paramsTypeFile, paramsType
}

func (s *Func) getParamsTypeByUrl(paramsUrl string, urlPars map[string]any, fsOpt *fsopts.Options) (paramsTypeFile string, paramsType common.ParamsType) {
	fsOpts := &fsopts.Options{
		RootPath: s.server.RootPath(),
		WorkPath: fsOpt.WorkPath,
	}

	fileUrl := paramsUrl + ".json"
	if (urlPars != nil) && (len(urlPars) > 0) {
		fileUrl = stringutils.ReplacePlaceholders(fileUrl, urlPars)
	}
	if !s.server.FsPkg().Exists(fileUrl, fsOpts) {
		return
	}
	paramsTypeFile = fileUrl
	bytes, err := s.server.ReadFile(fileUrl, fsOpts)
	if err != nil {
		panic(err)
	}
	if bytes != nil && len(bytes) > 0 {
		if err = jsonutils.Unmarshal(bytes, &paramsType); err != nil {
			panic(fmt.Sprintf(" loading %s  error: %s", fileUrl, err.Error()))
		}
		paramsType = paramsType
	}
	return paramsTypeFile, paramsType
}
