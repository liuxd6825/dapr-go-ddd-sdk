package entities

import (
	"github.com/jinzhu/copier"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/reflectutils"
)

func Copy[T any](src any, opts ...*copier.Option) T {
	target := reflectutils.NewInstance[T]()
	var option *copier.Option
	var err error
	for _, opt := range opts {
		option = opt
	}
	if option == nil {
		err = copier.Copy(target, src)
	} else {
		err = copier.CopyWithOption(target, src, *option)
	}
	if err != nil {
		panic(err)
	}
	return target
}
