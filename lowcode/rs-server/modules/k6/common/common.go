package common

import "github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"

type CommonPkg struct {
}

func NewCommonPkg() *CommonPkg {
	return &CommonPkg{}
}

func (c *CommonPkg) NewResult(data any, err error) *common.Result[any] {
	return &common.Result[any]{Data: data, Error: err}
}
