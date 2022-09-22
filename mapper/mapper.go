package mapper

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
)

//
// Mapper
// @Description: 进行struct属性复制，支持深度复制
// @param fromObj 来源
// @param toObj 目标
// @return error
//
func Mapper(fromObj, toObj interface{}) error {
	return types.Mapper(fromObj, toObj)
}

//
// MaskMapper
// @Description: 根据指定进行属性复制，不支持深度复制
// @param fromObj 来源
// @param toObj 目标
// @param mask 要复制属性列表
// @return error
//
func MaskMapper(fromObj, toObj interface{}, mask []string) error {
	return types.MaskMapper(fromObj, toObj, mask)
}

//
// NewMap
// @Description:
// @param formObj
// @return map[string]interface{}
// @return error
//
func NewMap(formObj interface{}) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	if err := MaskMapper(formObj, &res, nil); err != nil {
		return nil, err
	}
	return res, nil
}
