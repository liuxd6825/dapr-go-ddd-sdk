package db

import "fmt"

type IEntity interface {
	GetTenantId() string
	SetTenantId(v string)
	GetId() string
	SetId(v string)
}

type IDataSet interface {
	GetDataSet() map[string]any
	SetDataSet(val map[string]any)
}

type Entity struct {
	dataSet map[string]any
}

func NewEntity() IEntity {
	return &Entity{dataSet: make(map[string]any)}
}

func NewEntityWidth(dataSet map[string]any) IEntity {
	return &Entity{dataSet: dataSet}
}

func (e *Entity) SetDataSet(val map[string]any) {
	e.dataSet = val
}

func (e *Entity) GetDataSet() map[string]any {
	return e.dataSet
}

func (e *Entity) GetTenantId() string {
	return e.string("tenantId")
}

func (e *Entity) SetTenantId(val string) {
	e.setString("tenantId", val)
}

func (e *Entity) GetId() string {
	return e.string("id")
}

func (e *Entity) SetId(val string) {
	e.setString("id", val)
}

func (e *Entity) string(name string) string {
	val, ok := e.dataSet[name]
	if ok {
		return fmt.Sprintf("%s", val)
	}
	return ""
}

func (e *Entity) setString(name, val string) {
	e.dataSet[name] = val
}
