package ddd_mongodb

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rsql"
	"go.mongodb.org/mongo-driver/bson"
)

type filterItem struct {
	parent *filterItem
	name   string
	value  interface{}
	items  []*filterItem
}

func newFilterItem(parent *filterItem, name string) *filterItem {
	return &filterItem{
		name:   name,
		parent: parent,
		value:  nil,
		items:  make([]*filterItem, 0),
	}
}

func (i *filterItem) addChildItem(name string, value interface{}) *filterItem {
	newItem := newFilterItem(i, name)
	newItem.value = value
	i.items = append(i.items, newItem)
	return newItem
}

func (i *filterItem) getAndItem() {
}

func (f *filterItem) getValues(data map[string]interface{}) {
	if len(f.items) != 0 {
		array := make([]interface{}, len(f.items))
		for i, v := range f.items {
			item := ddd_utils.NewMap()
			item[v.name] = v.value
			array[i] = item
		}
		data[f.name] = array
	} else if f.value != nil {
		data[f.name] = f.value
	}
}

func (i *filterItem) setValue(name string, value interface{}) {
	i.name = name
	i.value = value
}

type MongoProcess struct {
	item    *filterItem
	current *filterItem
}

func NewMongoProcess() *MongoProcess {
	m := &MongoProcess{
		item: newFilterItem(nil, "$and"),
	}
	m.init()
	return m
}

func (m *MongoProcess) init() {
	m.current = m.item
}

func (m *MongoProcess) GetFilter(tenantId string) map[string]interface{} {
	data := make(map[string]interface{})
	m.item.getValues(data)
	m1, ok := data[""]
	if ok {
		d1 := m1.(map[string]interface{})
		d1["tenantId"] = tenantId
	} else if len(data) == 0 {
		data["tenantId"] = tenantId
	} else {
		m1, ok := data["$and"]
		d1, ok := m1.(map[string]interface{})
		if ok {
			d1["tenantId"] = tenantId
		}
		d2, ok := m1.([]interface{})
		if ok {
			item := ddd_utils.NewMap()
			item["tenantId"] = tenantId
			d2 := append(d2, item)
			data["$and"] = d2
		}
	}
	return data
}

func (m *MongoProcess) OnAndItem() {
	m.current.name = "$and"
}

func (m *MongoProcess) OnAndStart() {
	m.current = m.current.addChildItem("$and", nil)
}

func (m *MongoProcess) OnAndEnd() {
	m.current = m.current.parent
}

func (m *MongoProcess) OnOrItem() {
	m.current.name = "$or"
}

func (m *MongoProcess) OnOrStart() {
	m.current = m.current.addChildItem("$or", nil)
}

func (m *MongoProcess) OnOrEnd() {
	m.current = m.current.parent
}

func (m *MongoProcess) OnEquals(name string, value interface{}, rValue rsql.Value) {
	m.current.addChildItem(name, value)
}

func (m *MongoProcess) OnNotEquals(name string, value interface{}, rValue rsql.Value) {
	m.current.addChildItem(name, bson.D{{"!=", value}})
}

func (m *MongoProcess) OnLike(name string, value interface{}, rValue rsql.Value) {
	m.current.addChildItem(name, bson.D{{"$lt", value}})
}

func (m *MongoProcess) OnNotLike(name string, value interface{}, rValue rsql.Value) {
	m.current.addChildItem(name, bson.D{{"$lt", value}})
}

func (m *MongoProcess) OnGreaterThan(name string, value interface{}, rValue rsql.Value) {
	m.current.addChildItem(name, bson.D{{"$lt", value}})
}

func (m *MongoProcess) OnGreaterThanOrEquals(name string, value interface{}, rValue rsql.Value) {
	m.current.addChildItem(name, bson.D{{"$lt", value}})
}

func (m *MongoProcess) OnLessThan(name string, value interface{}, rValue rsql.Value) {
	m.current.addChildItem(name, bson.D{{"$lt", value}})
}

func (m *MongoProcess) OnLessThanOrEquals(name string, value interface{}, rValue rsql.Value) {
	m.current.addChildItem(name, bson.D{{"$lt", value}})
}

func (m *MongoProcess) OnIn(name string, value interface{}, rValue rsql.Value) {
	m.current.addChildItem(name, bson.D{{"$lt", value}})
}

func (m *MongoProcess) OnNotIn(name string, value interface{}, rValue rsql.Value) {
	m.current.addChildItem(name, bson.D{{"$lt", value}})
}
