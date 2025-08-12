package rsql_mongo

import (
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
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
	_id       = "id"
	tenant_id = "tenant_id"
)

func NewProcess(tenantId string) rsql.Process {
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
		return p.asTableName + "." + rsql.AsFieldName(name)
	}
	return rsql.AsFieldName(name)
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

func (p *Process) OnEquals(name string, value interface{}, rValue rsql.Value) {
	p.current.addChildItem(p.getFieldName(name), value)
}

func (p *Process) OnNotEquals(name string, value interface{}, rValue rsql.Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$ne", value}})
}

func (p *Process) OnLike(name string, value interface{}, rValue rsql.Value) {
	pattern := fmt.Sprintf("%s", value)
	pattern = strings.ReplaceAll(pattern, "*", "")

	p.current.addChildItem(p.getFieldName(name), primitive.Regex{Pattern: pattern, Options: "im"})
}

func (p *Process) OnNotLike(name string, value interface{}, rValue rsql.Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$lt", value}})
}

func (p *Process) OnGreaterThan(name string, value interface{}, rValue rsql.Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$gt", value}})
}

func (p *Process) OnGreaterThanOrEquals(name string, value interface{}, rValue rsql.Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$gte", value}})
}

func (p *Process) OnLessThan(name string, value interface{}, rValue rsql.Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$lt", value}})
}

func (p *Process) OnLessThanOrEquals(name string, value interface{}, rValue rsql.Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$lte", value}})
}

func (p *Process) OnIn(name string, value interface{}, rValue rsql.Value) {
	if listValue, ok := rValue.(*rsql.ListValue); ok {
		values := rsql.GetValueList(listValue)
		p.current.addChildItem(p.getFieldName(name), bson.M{"$in": values})
	}
}

func (p *Process) OnNotIn(name string, value interface{}, rValue rsql.Value) {
	if listValue, ok := rValue.(*rsql.ListValue); ok {
		values := rsql.GetValueList(listValue)
		p.current.addChildItem(p.getFieldName(name), bson.M{"$nin": values})
	}
}

func (p *Process) OnContains(name string, value interface{}, rValue rsql.Value) {
	val := fmt.Sprintf(".*%v.*", value)
	// "$regex": primitive.Regex{Pattern: ".*"+city+".*", Options: "i"}
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$regex", primitive.Regex{Pattern: val, Options: "i"}}})
}

func (p *Process) OnNotContains(name string, value interface{}, rValue rsql.Value) {
	val := fmt.Sprintf(".*%v.*", value)
	// "$regex": primitive.Regex{Pattern: ".*"+city+".*", Options: "i"}
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$not", primitive.Regex{Pattern: val, Options: "i"}}})
}

func (p *Process) OnFnProcess(expr rsql.Expression, fn *rsql.FuncValue) rsql.Value {
	if fn.Name == "sub" {
		iden, ok := expr.(rsql.GetIdentifier)
		if !ok {
			panic(errors.New("fn expression must implement GetIdentifier"))
		}

		foreignField := rsql.AsFieldName(fn.Args["field"])
		table := fn.Args["table"]
		tableAs := table + "_as"
		localField := iden.GetIdentifier().Val
		if localField == "id" {
			localField = "_id"
		}
		np := newProcess(p.tenantId)
		np.asTableName = tableAs
		np.current = p.current
		err := rsql.ParseProcess(fn.RSQL(), np)
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

func (p *Process) OnIsNull(name string, value interface{}, rValue rsql.Value) {
	p.current.addChildItem(p.getFieldName(name), nil)
}

func (p *Process) OnNotIsNull(name string, value interface{}, rValue rsql.Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$ne", nil}})
}

func (p *Process) OnStart(name string, value interface{}, rValue rsql.Value) {
	p.current.addChildItem(p.getFieldName(name), bson.D{{"$ne", nil}})
}

func (p *Process) OnEnd(name string, value interface{}, rValue rsql.Value) {
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

func (p *Process) getValue(value rsql.Value) interface{} {
	var v any
	var err error
	switch value.(type) {
	case *rsql.StringValue:
		sv, _ := value.(*rsql.StringValue)
		v = sv.Value
	case *rsql.IntegerValue:
		sv, _ := value.(*rsql.IntegerValue)
		v = sv.Value
	case *rsql.DateValue:
		sv, _ := value.(*rsql.DateValue)
		v, err = time.ParseInLocation(rsql.DateLayout, sv.Value, time.Local)
	case *rsql.DoubleValue:
		sv, _ := value.(*rsql.DoubleValue)
		v = sv.Value
	case *rsql.DateTimeValue:
		sv, _ := value.(*rsql.DateTimeValue)
		v, err = time.ParseInLocation(rsql.DateTimeLayout, sv.Value, time.Local)
	case *rsql.BooleanValue:
		sv, _ := value.(*rsql.BooleanValue)
		v = sv.Value
	case *rsql.ListValue:
		sv, _ := value.(*rsql.ListValue)
		v = rsql.GetValueList(sv)
	case *rsql.FuncValue:
		sv, _ := value.(*rsql.FuncValue)
		v = sv.Value
	default:
		v = value
	}
	if err != nil {
		panic(err)
	}
	return v
}
