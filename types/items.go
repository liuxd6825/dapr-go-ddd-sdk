package types

import (
	"context"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/mapper"
)

type Item interface {
	GetId() string
}

type Items[T Item] map[string]T

func (t *Items[T]) NewItem() T {
	return T{}
}

//
// Add
// @Description: 添加
// @param ctx    上下文
// @param data   更新数据
// @return error 错误
//
func (t *Items[T]) Add(ctx context.Context, id string, data interface{}) error {
	m := *t
	_, ok := m[id]
	if ok {
		return errors.New(fmt.Sprintf("新建 Id \"%s\" 已经存在", id))
	}
	newItem := t.NewItem()
	err := mapper.Mapper(data, newItem)
	if err != nil {
		m[id] = newItem
	}
	return err
}

//
// Update
// @Description:     更新
// @param ctx        上下文
// @param data       更新数据
// @param updateMask 更新字段项
// @return error
//
func (t Items[T]) Update(ctx context.Context, id string, data interface{}, updateMask []string) error {
	item, ok := t[id]
	if !ok {
		return nil
	}
	return mapper.MaskMapper(data, item, updateMask)
}

//
// Delete
// @Description: 删除明细
// @param ctx    上下文
// @param item   明细对象
// @return error 错误
//
func (t Items[T]) Delete(ctx context.Context, item T) error {
	delete(t, item.GetId())
	return nil
}

//
// DeleteById
// @Description: 按Id删除
// @param ctx    上下文
// @param id     Id主键
// @return error 错误
//
func (t Items[T]) DeleteById(ctx context.Context, id string) error {
	delete(t, id)
	return nil
}

//
// DeleteByIds
// @Description:  按id删除多个
// @receiver s
// @param ctx     上下文
// @param id      Id主键
// @return error  错误m
//
func (t Items[T]) DeleteByIds(ctx context.Context, ids ...string) error {
	if len(ids) > 0 {
		for _, id := range ids {
			if err := t.DeleteById(ctx, id); err != nil {
				return err
			}
		}
	}
	return nil
}
