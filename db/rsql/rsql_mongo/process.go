package rsql_mongo

import (
	"errors"
	"fmt"
	rsql2 "github.com/liuxd6825/dapr-go-ddd-sdk/db/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

type Process struct {
	item        *Item
	current     *Item
	errList     []string
	tenantId    string
	lookup      []*Lookup
	asTableName string
}

const (
	_id       = "_id"
	tenant_id = "tenant_id"
)

func NewProcess(tenantId string) rsql2.Process {
	return newProcess(tenantId)
}

func newProcess(tenantId string) *Process {
	m := &Process{
		item:     NewItem(nil, "$and"),
		errList:  make([]string, 0),
		tenantId: tenantId,
	}
	m.init()
	return m
}

func (p *Process) init() {
	p.current = p.item
}

func (p *Process) getFieldName(name string) string {
	if p.asTableName != "" {
		return p.asTableName + "." + rsql2.AsFieldName(name)
	}
	return rsql2.AsFieldName(name)
}

func (p *Process) OnAndItem() {
	p.current.name = "$and"
}

func (p *Process) OnAndStart() {
	p.current = p.current.addChildItem("$and", nil)
}

func (p *Process) OnAndEnd() {
	p.current = p.current.parent
}

func (p *Process) OnOrItem() {
	p.current.name = "$or"
}

func (p *Process) OnOrStart() {
	p.current = p.current.addChildItem("$or", nil)
}

func (p *Process) OnOrEnd() {
	p.current = p.current.parent
}

func (p *Process) OnEquals(name string, value interface{}, rValue rsql2.Value) {
	value = p.getValue(rValue)
	p.current.addChildItem(p.getFieldName(name), value)
}

func (p *Process) OnNotEquals(name string, value interface{}, rValue rsql2.Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$ne", p.getValue(rValue)}})
}

func (p *Process) OnLike(name string, value interface{}, rValue rsql2.Value) {
	value = p.getValue(rValue)
	pattern := fmt.Sprintf("%s", value)
	pattern = strings.ReplaceAll(pattern, "*", "")

	p.current.addChildItem(p.getFieldName(name), primitive.Regex{Pattern: pattern, Options: "im"})
}

func (p *Process) OnNotLike(name string, value interface{}, rValue rsql2.Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$lt", p.getValue(rValue)}})
}

func (p *Process) OnGreaterThan(name string, value interface{}, rValue rsql2.Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$gt", p.getValue(rValue)}})
}

func (p *Process) OnGreaterThanOrEquals(name string, value interface{}, rValue rsql2.Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$gte", p.getValue(rValue)}})
}

func (p *Process) OnLessThan(name string, value interface{}, rValue rsql2.Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$lt", p.getValue(rValue)}})
}

func (p *Process) OnLessThanOrEquals(name string, value interface{}, rValue rsql2.Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$lte", p.getValue(rValue)}})
}

func (p *Process) OnIn(name string, value interface{}, rValue rsql2.Value) {
	if listValue, ok := rValue.(*rsql2.ListValue); ok {
		values := rsql2.GetValueList(listValue)
		p.current.addChildItem(p.getFieldName(name), bson.M{"$in": values})
	}
}

func (p *Process) OnNotIn(name string, value interface{}, rValue rsql2.Value) {
	if listValue, ok := rValue.(*rsql2.ListValue); ok {
		values := rsql2.GetValueList(listValue)
		p.current.addChildItem(p.getFieldName(name), bson.M{"$nin": values})
	}
}

func (p *Process) OnContains(name string, value interface{}, rValue rsql2.Value) {
	val := fmt.Sprintf(".*%v.*", p.getValue(rValue))
	// "$regex": primitive.Regex{Pattern: ".*"+city+".*", Options: "i"}
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$regex", primitive.Regex{Pattern: val, Options: "i"}}})
}

func (p *Process) OnNotContains(name string, value interface{}, rValue rsql2.Value) {
	val := fmt.Sprintf(".*%v.*", p.getValue(rValue))
	// "$regex": primitive.Regex{Pattern: ".*"+city+".*", Options: "i"}
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$not", primitive.Regex{Pattern: val, Options: "i"}}})
}

func (p *Process) OnFnProcess(expr rsql2.Expression, fn *rsql2.FuncValue) rsql2.Value {
	if fn.Name == "sub" {
		iden, ok := expr.(rsql2.GetIdentifier)
		if !ok {
			panic(errors.New("fn expression must implement GetIdentifier"))
		}

		foreignField := rsql2.AsFieldName(fn.Args["field"])
		table := fn.Args["table"]
		tableAs := table + "_as"
		localField := iden.GetIdentifier().Val
		if localField == "id" {
			localField = "_id"
		}
		np := newProcess(p.tenantId)
		np.asTableName = tableAs
		np.current = p.current
		err := rsql2.ParseProcess(fn.RSQL(), np)
		if err != nil {
			panic(err)
		}

		p.lookup = append(p.lookup, &Lookup{
			From:         table,
			LocalField:   localField,
			ForeignField: foreignField,
			As:           tableAs,
		})
	}
	return nil
}

func (p *Process) OnIsNull(name string, value interface{}, rValue rsql2.Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$in", []interface{}{nil}}})
}

func (p *Process) OnNotIsNull(name string, value interface{}, rValue rsql2.Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$ne", nil}})
}

func (p *Process) OnStart(name string, value interface{}, rValue rsql2.Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$ne", nil}})
}

func (p *Process) OnEnd(name string, value interface{}, rValue rsql2.Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$ne", nil}})
}

func (p *Process) addChildItem(name string, value interface{}) {
	p.current.addChildItem(p.getFieldName(name), value)
}

func (p *Process) GetSQL() string {
	return ""
}

func (p *Process) GetFilter() any {
	filter := &Filter{}
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

	filter.Lookups = p.lookup
	filter.Match = match
	return filter
}

func (p *Process) addError(name string, err error) {
	if err != nil {
		msg := fmt.Sprintf("%v %v; ", name, err.Error())
		p.errList = append(p.errList, msg)
	}
}

func (p *Process) getValue(value rsql2.Value) interface{} {
	var v any
	var err error
	switch value.(type) {
	case *rsql2.StringValue:
		sv, _ := value.(*rsql2.StringValue)
		v = sv.Value
	case *rsql2.IntegerValue:
		sv, _ := value.(*rsql2.IntegerValue)
		v = sv.Value
	case *rsql2.DateValue:
		sv, _ := value.(*rsql2.DateValue)
		v, err = time.Parse(rsql2.DateLayout, sv.Value)
	case *rsql2.DoubleValue:
		sv, _ := value.(*rsql2.DoubleValue)
		v = sv.Value
	case *rsql2.DateTimeValue:
		sv, _ := value.(*rsql2.DateTimeValue)
		v, err = time.Parse(rsql2.DateTimeLayout, sv.Value)
	case *rsql2.BooleanValue:
		sv, _ := value.(*rsql2.BooleanValue)
		v = sv.Value
	case *rsql2.ListValue:
		sv, _ := value.(*rsql2.ListValue)
		v = rsql2.GetValueList(sv)
	case *rsql2.FuncValue:
		sv, _ := value.(*rsql2.FuncValue)
		v = sv.Value
	default:
		v = value
	}
	if err != nil {
		panic(err)
	}
	return v
}
