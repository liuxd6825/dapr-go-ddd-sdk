package rsql

import (
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

type mongoProcess struct {
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

func NewMongoProcess() Process {
	m := &mongoProcess{
		item:    newFilterItem(nil, "$and"),
		errList: make([]string, 0),
	}
	m.init()
	return m
}

func newFilterItem(parent *filterItem, name string) *filterItem {
	n := name
	if n == "id" {
		n = _id
	}
	return &filterItem{
		name:   n,
		parent: parent,
		value:  nil,
		items:  make([]*filterItem, 0),
	}
}

func (m *mongoProcess) init() {
	m.current = m.item
}

func (m *mongoProcess) OnAndItem() {
	m.current.name = "$and"
}

func (m *mongoProcess) OnAndStart() {
	m.current = m.current.addChildItem("$and", nil)
}

func (m *mongoProcess) OnAndEnd() {
	m.current = m.current.parent
}

func (m *mongoProcess) OnOrItem() {
	m.current.name = "$or"
}

func (m *mongoProcess) OnOrStart() {
	m.current = m.current.addChildItem("$or", nil)
}

func (m *mongoProcess) OnOrEnd() {
	m.current = m.current.parent
}

func (m *mongoProcess) OnEquals(name string, value interface{}, rValue Value) {
	value = getValue(rValue)
	m.current.addChildItem(AsFieldName(name), value)
}

func (m *mongoProcess) OnNotEquals(name string, value interface{}, rValue Value) {
	m.current.addChildItem(AsFieldName(name), bson.D{{"$ne", m.getValue(rValue)}})
}

func (m *mongoProcess) OnLike(name string, value interface{}, rValue Value) {
	value = getValue(rValue)
	pattern := fmt.Sprintf("%s", value)
	pattern = strings.ReplaceAll(pattern, "*", "")

	m.current.addChildItem(AsFieldName(name), primitive.Regex{Pattern: pattern, Options: "im"})
}

func (m *mongoProcess) OnNotLike(name string, value interface{}, rValue Value) {
	m.current.addChildItem(AsFieldName(name), bson.D{{"$lt", m.getValue(rValue)}})
}

func (m *mongoProcess) OnGreaterThan(name string, value interface{}, rValue Value) {
	m.current.addChildItem(AsFieldName(name), bson.D{{"$gt", m.getValue(rValue)}})
}

func (m *mongoProcess) OnGreaterThanOrEquals(name string, value interface{}, rValue Value) {
	m.current.addChildItem(AsFieldName(name), bson.D{{"$gte", m.getValue(rValue)}})
}

func (m *mongoProcess) OnLessThan(name string, value interface{}, rValue Value) {
	m.current.addChildItem(AsFieldName(name), bson.D{{"$lt", m.getValue(rValue)}})
}

func (m *mongoProcess) OnLessThanOrEquals(name string, value interface{}, rValue Value) {
	m.current.addChildItem(AsFieldName(name), bson.D{{"$lte", m.getValue(rValue)}})
}

func (m *mongoProcess) OnIn(name string, value interface{}, rValue Value) {
	listValue, _ := rValue.(ListValue)
	values := getValueList(listValue)
	m.current.addChildItem(AsFieldName(name), bson.M{"$in": values})
}

func (m *mongoProcess) OnNotIn(name string, value interface{}, rValue Value) {
	listValue, _ := rValue.(ListValue)
	values := getValueList(listValue)
	m.current.addChildItem(AsFieldName(name), bson.M{"$nin": values})
}

func (m *mongoProcess) OnContains(name string, value interface{}, rValue Value) {
	val := fmt.Sprintf(".*%v.*", m.getValue(rValue))
	// "$regex": primitive.Regex{Pattern: ".*"+city+".*", Options: "i"}
	m.current.addChildItem(AsFieldName(name), bson.D{{"$regex", primitive.Regex{Pattern: val, Options: "i"}}})
}

func (m *mongoProcess) OnNotContains(name string, value interface{}, rValue Value) {
	val := fmt.Sprintf(".*%v.*", m.getValue(rValue))
	// "$regex": primitive.Regex{Pattern: ".*"+city+".*", Options: "i"}
	m.current.addChildItem(AsFieldName(name), bson.D{{"$not", primitive.Regex{Pattern: val, Options: "i"}}})
}

func (m *mongoProcess) OnIsNull(name string, value interface{}, rValue Value) {
	m.current.addChildItem(AsFieldName(name), bson.D{{"$in", []interface{}{nil}}})
}

func (m *mongoProcess) OnNotIsNull(name string, value interface{}, rValue Value) {
	m.current.addChildItem(AsFieldName(name), bson.D{{"$ne", nil}})
}

func (m *mongoProcess) OnStart(name string, value interface{}, rValue Value) {
	m.current.addChildItem(AsFieldName(name), bson.D{{"$ne", nil}})
}

func (m *mongoProcess) OnEnd(name string, value interface{}, rValue Value) {
	m.current.addChildItem(AsFieldName(name), bson.D{{"$ne", nil}})
}

func (m *mongoProcess) GetSQL() string {
	return ""
}

func (m *mongoProcess) GetFilter(tenantId string) (map[string]any, error) {
	data := make(map[string]any)
	if len(m.errList) > 0 {
		msg := strings.Join(m.errList, " ")
		return nil, errors.New(msg)
	}

	m.item.getValues(data)
	m1, ok := data[""]
	if ok {
		d1 := m1.(map[string]any)
		d1[tenantId] = tenantId
	} else if len(data) == 0 {
		data[tenantId] = tenantId
	} else {
		m1, ok := data["$and"]
		d1, ok := m1.(map[string]interface{})
		if ok {
			d1[tenantId] = tenantId
		}
		d2, ok := m1.([]interface{})
		if ok {
			item := ddd_utils.NewMap()
			item[tenantId] = tenantId
			d2 := append(d2, item)
			data["$and"] = d2
		}
	}
	return data, nil
}

func (m *mongoProcess) addError(name string, err error) {
	if err != nil {
		msg := fmt.Sprintf("%v %v; ", name, err.Error())
		m.errList = append(m.errList, msg)
	}
}

func (m *mongoProcess) getValue(rValue Value) interface{} {
	v := getValue(rValue)
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

/*
func getValue(value Value) (interface{}, error) {
	var v interface{}
	var err error
	switch value.(type) {
	case StringValue:
		sv, _ := value.(StringValue)
		v = sv.Value
	case IntegerValue:
		sv, _ := value.(IntegerValue)
		v = sv.Value
	case DateValue:
		sv, _ := value.(DateValue)
		v, err = time.Parse(dateLayout, sv.Value)
	case DoubleValue:
		sv, _ := value.(DoubleValue)
		v = sv.Value
	case DateTimeValue:
		sv, _ := value.(DateTimeValue)
		v, err = time.Parse(dateTimeLayout, sv.Value)
	case BooleanValue:
		sv, _ := value.(BooleanValue)
		v = sv.Value
	case ListValue:
		sv, _ := value.(ListValue)
		v, err = getValueList(sv)
	default:
		v = value
	}
	return v, err
}

func getValueList(listValue ListValue) ([]interface{}, error) {
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
*/
