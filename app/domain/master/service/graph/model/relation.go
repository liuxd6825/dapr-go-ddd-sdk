package model

import (
	"errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
)

type Relation struct {
	Id          string `json:"id" gorm:"column:id"`
	CaseId      string `json:"caseId" gorm:"column:case_id;nodeLabel:true;nodeLabelFormat:case_%"`
	RelType     string `json:"relType" gorm:"column:rel_type;relType:true"`
	StartId     string `json:"startId" gorm:"column:start_id;relStartId:true"`
	EndId       string `json:"endId" gorm:"column:end_id;relEndId:true"`
	TenantId    string `json:"tenantId" gorm:"column:tenant_id;nodeLabel:true;nodeLabelFormat:tenant_%"`
	TableName   string `json:"tableName" gorm:"column:table_name"`
	Description string `json:"description" gorm:"column:description"`
}

func NewRelation(dbSch *dbschema.DBSchema, data map[string]any) *Relation {
	if dbSch == nil {
		panic(errors.New("dbSch is nil"))
	}
	tenantId, _ := maputils.GetString(data, "tenantId", "")
	caseId, _ := maputils.GetString(data, "caseId", "")
	id, _ := maputils.GetString(data, "id", "")
	relStartId, _ := getRelStartId(dbSch, data)
	relEndId, _ := getRelEndId(dbSch, data)
	relType, _ := getRelType(dbSch, data)
	return &Relation{
		Id:          id,
		CaseId:      caseId,
		TenantId:    tenantId,
		RelType:     relType,
		StartId:     relStartId,
		EndId:       relEndId,
		TableName:   dbSch.TableName,
		Description: description(data),
	}
}

func getRelStartId(dbSch *dbschema.DBSchema, data map[string]any) (string, error) {
	field := dbSch.GetRelStartIdField()
	if field == nil {
		return "", errors.New("start id not found")
	}
	val, err := maputils.GetString(data, field.Name, "")
	if err != nil {
		return "", err
	}
	return val, nil
}

func getRelEndId(dbSch *dbschema.DBSchema, data map[string]any) (string, error) {
	field := dbSch.GetRelEndIdField()
	if field == nil {
		return "", errors.New("start id not found")
	}
	val, err := maputils.GetString(data, field.Name, "")
	if err != nil {
		return "", err
	}
	return val, nil
}

func getRelType(dbSch *dbschema.DBSchema, data map[string]any) (string, error) {
	field := dbSch.GetRelTypeField()
	if field == nil {
		return "", errors.New("start id not found")
	}
	val, err := maputils.GetString(data, field.Name, "")
	if err != nil {
		return "", err
	}
	return val, nil
}
