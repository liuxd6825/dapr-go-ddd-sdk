package runtime

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/runtime/transform"
)

func RecoverError(e error, recover any) error {
	var err error
	if e != nil {
		err = e
	} else if recover != nil {
		if ve, ok := recover.(error); ok {
			err = ve
		} else if obj, ok := recover.(*goja.Object); ok {
			eObj := obj.Export()
			err = fmt.Errorf("unknown error %s", eObj)
		} else {
			err = fmt.Errorf("unknown error %s", recover)
		}
	}
	return err
}

type TransformType = int

const (
	TransformTypeTypeScript TransformType = iota
	TransformTypeES6
)

// TransformCode
//
//	@Description: 将typescript代码转换为js
//	@param tsCode
//	@return string
//	@return error
func TransformCode(tsCode string, fileName string, transType TransformType) ([]byte, []*transform.FuncParam, error) {
	var codeBytes []byte
	var err error
	var params []*transform.FuncParam
	switch transType {
	case TransformTypeTypeScript:
		codeBytes, params, err = transform.TransformFromTypeScript(tsCode, fileName)
	case TransformTypeES6:
		codeBytes, params, err = transform.TransformFromEs6(tsCode, fileName)
	}
	return codeBytes, params, err
}
