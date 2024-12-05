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

// TransformTSCodeToJS
//
//	@Description: 将typescript代码转换为js
//	@param tsCode
//	@return string
//	@return error
func TransformTSCodeToJS(tsCode string) ([]byte, error) {
	jscode, err := transform.Transform(tsCode)
	return jscode, err
}
