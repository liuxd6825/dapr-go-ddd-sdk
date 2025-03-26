package runtime

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsopts"
	"github.com/spf13/afero"
	"io"
	"os"
	"testing"
)

var code = `
 	feign = require("/js/human-feign.js");
`

func TestRuntime_RunString(t *testing.T) {
	reader := &FileReader{}
	vm := NewRuntime(reader)
	resData, err := vm.Require("/js/human-feign1.js")
	if err != nil {
		t.Fatal(err)
	}
	vm1 := vm.GetVM()
	vm.PrintValue(resData)

	defVal := resData.ToObject(vm1).Get("default")
	vm.PrintValue(defVal)

	// 获取原型方法
	if fValue := resData.Export(); fValue != nil {
		t.Log(fValue)
	}
	t.Log(resData)
}

type FileReader struct {
}

func (r *FileReader) ReadText(filename string, opts ...*fsopts.Options) (string, error) {
	data, err := r.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (r *FileReader) ReadFile(filename string, opts ...*fsopts.Options) ([]byte, error) {
	root, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	file, err := afero.NewOsFs().Open(root + filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}
