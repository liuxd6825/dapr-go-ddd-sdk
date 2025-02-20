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

func (p *mongoProcess) init() {
	p.current = p.item
}

func (p *mongoProcess) OnFnProcess(fn *FuncValue) Value {
	return fn
}

func (p *mongoProcess) OnAndItem() {
	p.current.name = "$and"
}

func (p *mongoProcess) OnAndStart() {
	p.current = p.current.addChildItem("$and", nil)
}

func (p *mongoProcess) OnAndEnd() {
	p.current = p.current.parent
}

func (p *mongoProcess) OnOrItem() {
	p.current.name = "$or"
}

func (p *mongoProcess) OnOrStart() {
	p.current = p.current.addChildItem("$or", nil)
}

func (p *mongoProcess) OnOrEnd() {
	p.current = p.current.parent
}

func (p *mongoProcess) OnEquals(name string, value interface{}, rValue Value) {
	value = getValue(rValue)
	p.current.addChildItem(AsFieldName(name), value)
}

func (p *mongoProcess) OnNotEquals(name string, value interface{}, rValue Value) {
	p.current.addChildItem(AsFieldName(name), bson.D{{"$ne", p.getValue(rValue)}})
}

func (p *mongoProcess) OnLike(name string, value interface{}, rValue Value) {
	value = getValue(rValue)
	pattern := fmt.Sprintf("%s", value)
	pattern = strings.ReplaceAll(pattern, "*", "")

	p.current.addChildItem(AsFieldName(name), primitive.Regex{Pattern: pattern, Options: "im"})
}

func (p *mongoProcess) OnNotLike(name string, value interface{}, rValue Value) {
	p.current.addChildItem(AsFieldName(name), bson.D{{"$lt", p.getValue(rValue)}})
}

func (p *mongoProcess) OnGreaterThan(name string, value interface{}, rValue Value) {
	p.current.addChildItem(AsFieldName(name), bson.D{{"$gt", p.getValue(rValue)}})
}

func (p *mongoProcess) OnGreaterThanOrEquals(name string, value interface{}, rValue Value) {
	p.current.addChildItem(AsFieldName(name), bson.D{{"$gte", p.getValue(rValue)}})
}

func (p *mongoProcess) OnLessThan(name string, value interface{}, rValue Value) {
	p.current.addChildItem(AsFieldName(name), bson.D{{"$lt", p.getValue(rValue)}})
}

func (p *mongoProcess) OnLessThanOrEquals(name string, value interface{}, rValue Value) {
	p.current.addChildItem(AsFieldName(name), bson.D{{"$lte", p.getValue(rValue)}})
}

func (p *mongoProcess) OnIn(name string, value interface{}, rValue Value) {
	listValue, _ := rValue.(*ListValue)
	values := getValueList(listValue)
	p.current.addChildItem(AsFieldName(name), bson.M{"$in": values})
}

func (p *mongoProcess) OnNotIn(name string, value interface{}, rValue Value) {
	listValue, _ := rValue.(*ListValue)
	values := getValueList(listValue)
	p.current.addChildItem(AsFieldName(name), bson.M{"$nin": values})
}

func (p *mongoProcess) OnContains(name string, value interface{}, rValue Value) {
	val := fmt.Sprintf(".*%v.*", p.getValue(rValue))
	// "$regex": primitive.Regex{Pattern: ".*"+city+".*", Options: "i"}
	p.current.addChildItem(AsFieldName(name), bson.D{{"$regex", primitive.Regex{Pattern: val, Options: "i"}}})
}

func (p *mongoProcess) OnNotContains(name string, value interface{}, rValue Value) {
	val := fmt.Sprintf(".*%v.*", p.getValue(rValue))
	// "$regex": primitive.Regex{Pattern: ".*"+city+".*", Options: "i"}
	p.current.addChildItem(AsFieldName(name), bson.D{{"$not", primitive.Regex{Pattern: val, Options: "i"}}})
}

func (p *mongoProcess) OnIsNull(name string, value interface{}, rValue Value) {
	p.current.addChildItem(AsFieldName(name), bson.D{{"$in", []interface{}{nil}}})
}

func (p *mongoProcess) OnNotIsNull(name string, value interface{}, rValue Value) {
	p.current.addChildItem(AsFieldName(name), bson.D{{"$ne", nil}})
}

func (p *mongoProcess) OnStart(name string, value interface{}, rValue Value) {
	p.current.addChildItem(AsFieldName(name), bson.D{{"$ne", nil}})
}

func (p *mongoProcess) OnEnd(name string, value interface{}, rValue Value) {
	p.current.addChildItem(AsFieldName(name), bson.D{{"$ne", nil}})
}

func (p *mongoProcess) OnInSubTable(name string, value interface{}, rValue Value) {

}

func (p *mongoProcess) GetSQL() string {
	return ""
}

func (p *mongoProcess) GetFilter(tenantId string) (map[string]any, error) {
	data := make(map[string]any)
	if len(p.errList) > 0 {
		msg := strings.Join(p.errList, " ")
		return nil, errors.New(msg)
	}

	p.item.getValues(data)
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

func (p *mongoProcess) addError(name string, err error) {
	if err != nil {
		msg := fmt.Sprintf("%v %v; ", name, err.Error())
		p.errList = append(p.errList, msg)
	}
}

func (p *mongoProcess) getValue(rValue Value) interface{} {
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
