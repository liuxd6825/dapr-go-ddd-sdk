package types

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type Item interface {
	GetId() string
}

type Items[T Item] struct {
	null    T
	items   map[string]T
	newFunc func() interface{}
}

func NewItems[T Item](newFunc func() interface{}) Items[T] {
	res := Items[T]{items: make(map[string]T)}
	err := res.Init(newFunc)
	if err != nil {
		panic(err)
	}
	return res
}

func (t *Items[T]) Init(newFunc func() interface{}) error {
	t.items = make(map[string]T)
	t.newFunc = newFunc
	return nil
}

func (t *Items[T]) NewItem() (T, error) {
	item, ok := t.newFunc().(T)
	if !ok {
		return t.null, fmt.Errorf("types.Items.NewItem() error ")
	}
	return item, nil
}

//
// AddMapper
// @Description: 添加
// @param ctx    上下文
// @param data   更新数据
// @return error 错误
//
func (t *Items[T]) AddMapper(ctx context.Context, id string, data interface{}) (T, error) {
	_, ok := t.items[id]
	if ok {
		return t.null, errors.New(fmt.Sprintf("新建 Id \"%s\" 已经存在", id))
	}
	newItem, err := t.NewItem()
	if err != nil {
		return t.null, err
	}
	err = Mapper(data, newItem)
	if err == nil {
		t.items[id] = newItem
	}
	return newItem, err
}

//
// AddItem
// @Description: 添加
// @param ctx    上下文
// @param data   更新数据
// @return error 错误
//
func (t *Items[T]) AddItem(ctx context.Context, item T) error {
	id := item.GetId()
	_, ok := t.items[id]
	if ok {
		return errors.New(fmt.Sprintf("新建 Id \"%s\" 已经存在", id))
	}
	t.items[id] = item
	return nil
}

//
// UpdateMapper
// @Description:     更新
// @param ctx        上下文
// @param data       更新数据
// @param updateMask 更新字段项
// @return error
//
func (t *Items[T]) UpdateMapper(ctx context.Context, id string, data interface{}, updateMask []string) (T, bool, error) {
	item, ok := t.items[id]
	if !ok {
		return item, ok, fmt.Errorf("types.Items.UpdateMapper() id %s ", id)
	}
	err := MaskMapper(data, item, updateMask)
	return item, ok, err
}

//
// UpdateItem
// @Description: 更新
// @param ctx    上下文
// @param data   更新数据
// @return error 错误
//
func (t *Items[T]) UpdateItem(ctx context.Context, item T) error {
	id := item.GetId()
	_, ok := t.items[id]
	if !ok {
		return errors.New(fmt.Sprintf("types.Items.UpdateItem()  Id \"%s\" 不存在", id))
	}
	t.items[id] = item
	return nil
}

//
// Delete
// @Description: 删除明细
// @param ctx    上下文
// @param item   明细对象
// @return error 错误
//
func (t *Items[T]) Delete(ctx context.Context, item T) error {
	delete(t.items, item.GetId())
	return nil
}

//
// DeleteById
// @Description: 按Id删除
// @param ctx    上下文
// @param id     Id主键
// @return error 错误
//
func (t *Items[T]) DeleteById(ctx context.Context, id string) error {
	delete(t.items, id)
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
func (t *Items[T]) DeleteByIds(ctx context.Context, ids ...string) error {
	if len(ids) > 0 {
		for _, id := range ids {
			if err := t.DeleteById(ctx, id); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Items[T]) ContainsId(id string) bool {
	_, ok := t.items[id]
	return ok
}

func (t *Items[T]) Get(id string) (T, bool) {
	item, ok := t.items[id]
	return item, ok
}

func (t *Items[T]) MapData() map[string]T {
	return t.items
}

func (t *Items[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.items)
}
