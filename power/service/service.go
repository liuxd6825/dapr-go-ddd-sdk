package service

import (
	"fmt"
	"github.com/dop251/goja"
	_ "github.com/dop251/goja_nodejs/console"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Service struct {
	srcFile string
	vm      *goja.Runtime
	opts    []ServiceOption
}

type ServiceOption func(s *Service) error

type ResultCode struct {
	Imports []string `yaml:"imports"`
	Code    []byte   `yaml:"code"`
}

func NewService(srcFile string, opts ...ServiceOption) *Service {
	return &Service{srcFile: srcFile, opts: opts}
}

func (s *Service) Run() error {
	s.vm = goja.New()
	loader := NewLoader(s.vm)
	for _, opt := range s.opts {
		if err := opt(s); err != nil {
			return err
		}
	}
	if err := loader.Load(s.srcFile); err != nil {
		fmt.Println("执行 JavaScript 代码时出错：", err)
		return err
	}
	return nil
}

func (s *Service) GetValue(name string) goja.Value {
	value := s.vm.Get(name)
	return value
}

func (s *Service) SetValue(name string, val any) error {
	if err := s.vm.Set(name, val); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (s *Service) Restart() error {
	return s.Run()
}

func (s *Service) GetFunc(funName string) (goja.Callable, bool) {
	fun := s.vm.Get(funName)
	return goja.AssertFunction(fun)
}

func (s *Service) readCode(filename string) (*ResultCode, error) {
	res := &ResultCode{}
	fileName, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}
	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	code := string(fileBytes)
	lines := strings.Split(code, "\n")
	for i, line := range lines {
		// 删除掉代码中: import {*} from "./types/*"
		if strings.HasPrefix(line, "import ") {
			iFrom := strings.Index(line, "from")
			if iFrom > 0 {
				if strings.Contains(line, "./types/") {
					lines[i] = "// " + line
				} else {
					src := line[iFrom+5:]
					src = strings.ReplaceAll(src, "\"", "")
					src = strings.ReplaceAll(src, ";", "")
					res.Imports = append(res.Imports, src)
					lines[i] = "// " + line
				}
			}
		}
		if i > 50 {
			break
		}
	}
	res.Code = []byte(strings.Join(lines, "\n"))
	return res, err
}
