package ddd_mongodb

import (
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rsql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

type MongoProcess struct {
	item    *filterItem
	current *filterItem
	errList []string
}

type filterItem struct {
	parent *filterItem
	name   string
	value  interface{}
	items  []*filterItem
}

const (
	dateTimeLayout = "2006-01-02T15:04:05"
	dateLayout     = "2006-01-02"
)

func NewMongoProcess() *MongoProcess {
	m := &MongoProcess{
		item:    newFilterItem(nil, "$and"),
		errList: make([]string, 0),
	}
	m.init()
	return m
}

func newFilterItem(parent *filterItem, name string) *filterItem {
	n := name
	if n == "id" {
		n = "_id"
	}
	return &filterItem{
		name:   n,
		parent: parent,
		value:  nil,
		items:  make([]*filterItem, 0),
	}
}

func (m *MongoProcess) init() {
	m.current = m.item
}

func (m *MongoProcess) GetFilter(tenantId string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	if len(m.errList) > 0 {
		msg := strings.Join(m.errList, " ")
		return nil, errors.New(msg)
	}

	m.item.getValues(data)
	m1, ok := data[""]
	if ok {
		d1 := m1.(map[string]interface{})
		d1[ConstTenantIdField] = tenantId
	} else if len(data) == 0 {
		data[ConstTenantIdField] = tenantId
	} else {
		m1, ok := data["$and"]
		d1, ok := m1.(map[string]interface{})
		if ok {
			d1[ConstTenantIdField] = tenantId
		}
		d2, ok := m1.([]interface{})
		if ok {
			item := ddd_utils.NewMap()
			item[ConstTenantIdField] = tenantId
			d2 := append(d2, item)
			data["$and"] = d2
		}
	}
	return data, nil
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
	value, err := getValue(rValue)
	if err != nil {
		m.addError(name, err)
	}
	m.current.addChildItem(AsFieldName(name), value)
}

func (m *MongoProcess) OnNotEquals(name string, value interface{}, rValue rsql.Value) {
	m.current.addChildItem(AsFieldName(name), bson.D{{"$ne", m.getValue(rValue)}})
}

func (m *MongoProcess) OnLike(name string, value interface{}, rValue rsql.Value) {
	value, err := getValue(rValue)
	if err != nil {
		m.addError(name, err)
	}
	pattern := fmt.Sprintf("%s", value)
	m.current.addChildItem(AsFieldName(name), primitive.Regex{Pattern: pattern, Options: "im"})
}

func (m *MongoProcess) OnNotLike(name string, value interface{}, rValue rsql.Value) {
	m.current.addChildItem(AsFieldName(name), bson.D{{"$lt", m.getValue(rValue)}})
}

func (m *MongoProcess) OnGreaterThan(name string, value interface{}, rValue rsql.Value) {
	m.current.addChildItem(AsFieldName(name), bson.D{{"$gt", m.getValue(rValue)}})
}

func (m *MongoProcess) OnGreaterThanOrEquals(name string, value interface{}, rValue rsql.Value) {
	m.current.addChildItem(AsFieldName(name), bson.D{{"$gte", m.getValue(rValue)}})
}

func (m *MongoProcess) OnLessThan(name string, value interface{}, rValue rsql.Value) {
	m.current.addChildItem(AsFieldName(name), bson.D{{"$lt", m.getValue(rValue)}})
}

func (m *MongoProcess) OnLessThanOrEquals(name string, value interface{}, rValue rsql.Value) {
	m.current.addChildItem(AsFieldName(name), bson.D{{"$lte", m.getValue(rValue)}})
}

func (m *MongoProcess) OnIn(name string, value interface{}, rValue rsql.Value) {
	listValue, _ := rValue.(rsql.ListValue)
	values, err := getValueList(listValue)
	if err != nil {
		m.addError(name, err)
	}
	m.current.addChildItem(AsFieldName(name), bson.M{"$in": values})
}

func (m *MongoProcess) OnNotIn(name string, value interface{}, rValue rsql.Value) {
	listValue, _ := rValue.(rsql.ListValue)
	values, err := getValueList(listValue)
	if err != nil {
		m.addError(name, err)
	}
	m.current.addChildItem(AsFieldName(name), bson.M{"$nin": values})
}

func (m *MongoProcess) addError(name string, err error) {
	if err != nil {
		msg := fmt.Sprintf("%v %v; ", name, err.Error())
		m.errList = append(m.errList, msg)
	}
}

func (m *MongoProcess) getValue(rValue rsql.Value) interface{} {
	v, err := getValue(rValue)
	if err != nil {
		m.addError(rValue.ValueName(), err)
	}
	return v
}

func (i *filterItem) addChildItem(name string, value interface{}) *filterItem {
	newItem := newFilterItem(i, name)
	newItem.value = value
	i.items = append(i.items, newItem)
	return newItem
}

func (i *filterItem) getAndItem() {
}

func (i *filterItem) getValues(data map[string]interface{}) {
	if len(i.items) != 0 {
		array := make([]interface{}, len(i.items))
		for i, v := range i.items {
			item := ddd_utils.NewMap()
			item[v.name] = v.value
			if len(v.items) > 0 {
				m := ddd_utils.NewMap()
				v.getValues(m)
				item[v.name] = m[v.name]
			}
			array[i] = item
		}
		data[i.name] = array
	} else if i.value != nil {
		data[i.name] = i.value
	}
}

func (i *filterItem) setValue(name string, value interface{}) {
	i.name = name
	i.value = value
}

func getValue(value rsql.Value) (interface{}, error) {
	var v interface{}
	var err error
	switch value.(type) {
	case rsql.StringValue:
		sv, _ := value.(rsql.StringValue)
		v = sv.Value
	case rsql.IntegerValue:
		sv, _ := value.(rsql.IntegerValue)
		v = sv.Value
	case rsql.DateValue:
		sv, _ := value.(rsql.DateValue)
		v, err = time.Parse(dateLayout, sv.Value)
	case rsql.DoubleValue:
		sv, _ := value.(rsql.DoubleValue)
		v = sv.Value
	case rsql.DateTimeValue:
		sv, _ := value.(rsql.DateTimeValue)
		v, err = time.Parse(dateTimeLayout, sv.Value)
	case rsql.BooleanValue:
		sv, _ := value.(rsql.BooleanValue)
		v = sv.Value
	case rsql.ListValue:
		sv, _ := value.(rsql.ListValue)
		v, err = getValueList(sv)
	default:
		v = value
	}
	return v, err
}

func getValueList(listValue rsql.ListValue) ([]interface{}, error) {
	list := make([]interface{}, 0)
	for _, item := range listValue.Value {
		v, err := getValue(item)
		if err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	return list, nil
}
