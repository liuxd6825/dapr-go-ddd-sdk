package service

import (
	"github.com/dop251/goja"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Loader struct {
	src string
	vm  *goja.Runtime
}

func NewLoader(vm *goja.Runtime) *Loader {
	return &Loader{vm: vm}
}

func (l *Loader) Load(src string) error {
	result, err := l.readFile(src)
	if err != nil {
		return err
	}
	for _, src := range result.Imports {
		src = strings.ReplaceAll(src, "./", "./services/")
		fileName, err := filepath.Abs(src + ".js")
		if err != nil {
			return err
		}
		err = l.Load(fileName)
		if err != nil {
			return err
		}
	}
	_, err = l.vm.RunScript(src, string(result.Code))
	return err
}

func (l *Loader) readFile(filename string) (*ResultCode, error) {
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
	code = strings.Join(lines, "\n")
	code = strings.ReplaceAll(code, "export ", "")
	res.Code = []byte(code)
	return res, err
}
