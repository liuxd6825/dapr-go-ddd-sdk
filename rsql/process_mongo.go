package rsql

import (
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

type mongoLookup struct {
	from         string
	localField   string
	foreignField string
	as           string
}
type mongoProcess struct {
	item        *filterItem
	current     *filterItem
	errList     []string
	tenantId    string
	lookup      []*mongoLookup
	asTableName string
}

type filterItem struct {
	parent *filterItem
	name   string
	value  interface{}
	items  []*filterItem
}

func NewMongoProcess(tenantId string) Process {
	return newMongoProcess(tenantId)
}

func newMongoProcess(tenantId string) *mongoProcess {
	m := &mongoProcess{
		item:     newFilterItem(nil, "$and"),
		errList:  make([]string, 0),
		tenantId: tenantId,
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

func (p *mongoProcess) getFieldName(name string) string {
	if p.asTableName != "" {
		return p.asTableName + "." + AsFieldName(name)
	}
	return AsFieldName(name)
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
	p.current.addChildItem(p.getFieldName(name), value)
}

func (p *mongoProcess) OnNotEquals(name string, value interface{}, rValue Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$ne", p.getValue(rValue)}})
}

func (p *mongoProcess) OnLike(name string, value interface{}, rValue Value) {
	value = getValue(rValue)
	pattern := fmt.Sprintf("%s", value)
	pattern = strings.ReplaceAll(pattern, "*", "")

	p.current.addChildItem(p.getFieldName(name), primitive.Regex{Pattern: pattern, Options: "im"})
}

func (p *mongoProcess) OnNotLike(name string, value interface{}, rValue Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$lt", p.getValue(rValue)}})
}

func (p *mongoProcess) OnGreaterThan(name string, value interface{}, rValue Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$gt", p.getValue(rValue)}})
}

func (p *mongoProcess) OnGreaterThanOrEquals(name string, value interface{}, rValue Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$gte", p.getValue(rValue)}})
}

func (p *mongoProcess) OnLessThan(name string, value interface{}, rValue Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$lt", p.getValue(rValue)}})
}

func (p *mongoProcess) OnLessThanOrEquals(name string, value interface{}, rValue Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$lte", p.getValue(rValue)}})
}

func (p *mongoProcess) OnIn(name string, value interface{}, rValue Value) {
	if listValue, ok := rValue.(*ListValue); ok {
		values := getValueList(listValue)
		p.current.addChildItem(p.getFieldName(name), bson.M{"$in": values})
	}
}

func (p *mongoProcess) OnNotIn(name string, value interface{}, rValue Value) {
	if listValue, ok := rValue.(*ListValue); ok {
		values := getValueList(listValue)
		p.current.addChildItem(p.getFieldName(name), bson.M{"$nin": values})
	}
}

func (p *mongoProcess) OnContains(name string, value interface{}, rValue Value) {
	val := fmt.Sprintf(".*%v.*", p.getValue(rValue))
	// "$regex": primitive.Regex{Pattern: ".*"+city+".*", Options: "i"}
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$regex", primitive.Regex{Pattern: val, Options: "i"}}})
}

func (p *mongoProcess) OnNotContains(name string, value interface{}, rValue Value) {
	val := fmt.Sprintf(".*%v.*", p.getValue(rValue))
	// "$regex": primitive.Regex{Pattern: ".*"+city+".*", Options: "i"}
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$not", primitive.Regex{Pattern: val, Options: "i"}}})
}

func (p *mongoProcess) OnFnProcess(expr Expression, fn *FuncValue) Value {
	switch fn.Name {
	case "sub":
		iden, ok := expr.(GetIdentifier)
		if !ok {
			panic(errors.New("fn expression must implement GetIdentifier"))
		}

		foreignField := AsFieldName(fn.Args["field"])
		table := fn.Args["table"]
		tableAs := table + "_as"
		localField := iden.GetIdentifier().Val
		if localField == "id" {
			localField = "_id"
		}
		np := newMongoProcess(p.tenantId)
		np.asTableName = tableAs
		np.current = p.current
		err := ParseProcess(fn.RSQL(), np)
		if err != nil {
			panic(err)
		}

		p.lookup = append(p.lookup, &mongoLookup{
			from:         table,
			localField:   localField,
			foreignField: foreignField,
			as:           tableAs,
		})

		//filter := np.GetFilter()
		//p.current.addChildItem("$match", filter)
		return nil
	default:
		return nil
	}
	return fn
}

func (p *mongoProcess) OnIsNull(name string, value interface{}, rValue Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$in", []interface{}{nil}}})
}

func (p *mongoProcess) OnNotIsNull(name string, value interface{}, rValue Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$ne", nil}})
}

func (p *mongoProcess) OnStart(name string, value interface{}, rValue Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$ne", nil}})
}

func (p *mongoProcess) OnEnd(name string, value interface{}, rValue Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$ne", nil}})
}

func (p *mongoProcess) addChildItem(name string, value interface{}) {
	p.current.addChildItem(p.getFieldName(name), value)
}
func (p *mongoProcess) GetSQL() string {
	return ""
}

func (p *mongoProcess) GetFilter() any {

	if len(p.errList) > 0 {
		msg := strings.Join(p.errList, " ")
		panic(errors.New(msg))
	}
	tenantIdField := p.getFieldName(tenant_id)
	match := make(map[string]any)
	p.item.getValues(match)
	m1, ok := match[""]
	if ok {
		d1 := m1.(map[string]any)
		d1[tenantIdField] = p.tenantId
	} else if len(match) == 0 {
		match[tenantIdField] = p.tenantId
	} else {
		m1, ok := match["$and"]
		d1, ok := m1.(map[string]any)
		if ok {
			d1[tenantIdField] = p.tenantId
		}
		d2, ok := m1.([]any)
		if ok {
			item := ddd_utils.NewMap()
			item[tenantIdField] = p.tenantId
			d2 := append(d2, item)
			match["$and"] = d2
		}
	}
	if len(p.lookup) == 0 {
		return match
	}

	var list []map[string]any
	// $lookup必须放在$match前面
	for _, v := range p.lookup {
		list = append(list, map[string]any{
			"$lookup": bson.M{
				"from":         v.from,
				"localField":   v.localField,
				"foreignField": v.foreignField,
				"as":           v.as,
			},
		})
		list = append(list, map[string]any{
			"$unwind": "$" + v.as,
		})
	}
	list = append(list, bson.M{"$match": match})
	return list
}

func (p *mongoProcess) addError(name string, err error) {
	if err != nil {
		msg := fmt.Sprintf("%v %v; ", name, err.Error())
		p.errList = append(p.errList, msg)
	}
}

func (p *mongoProcess) getValue(value Value) interface{} {
	var v any
	var err error
	switch value.(type) {
	case *StringValue:
		sv, _ := value.(*StringValue)
		v = sv.Value
	case *IntegerValue:
		sv, _ := value.(*IntegerValue)
		v = sv.Value
	case *DateValue:
		sv, _ := value.(*DateValue)
		v, err = time.Parse(dateLayout, sv.Value)
	case *DoubleValue:
		sv, _ := value.(*DoubleValue)
		v = sv.Value
	case *DateTimeValue:
		sv, _ := value.(*DateTimeValue)
		v, err = time.Parse(dateTimeLayout, sv.Value)
	case *BooleanValue:
		sv, _ := value.(*BooleanValue)
		v = sv.Value
	case *ListValue:
		sv, _ := value.(*ListValue)
		v = getValueList(sv)
	case *FuncValue:
		sv, _ := value.(*FuncValue)
		v = sv.Value
	default:
		v = value
	}
	if err != nil {
		panic(err)
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
