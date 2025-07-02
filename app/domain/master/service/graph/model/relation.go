package model

import (
	"errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/graph/utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
)

type Relation struct {
	Id          string   `json:"id" gorm:"column:id"`
	CaseId      string   `json:"caseId" gorm:"column:case_id;nodeLabel:true;nodeLabelFormat:case_%"`
	TenantId    string   `json:"tenantId" gorm:"column:tenant_id;nodeLabel:true;nodeLabelFormat:tenant_%"`
	Source      string   `json:"source" gorm:"column:source;relStartId:true"`
	Target      string   `json:"target" gorm:"column:target;relEndId:true"`
	SourceIds   string   `json:"sourceIds" gorm:"column:source_ids"`
	SourceType  string   `json:"sourceType" gorm:"column:source_type"`
	RelType     string   `json:"relType" gorm:"column:rel_type;relType:true"`
	Keywords    []string `json:"keywords" gorm:"column:keywords;type:text;serializer:json"`
	Description string   `json:"description" gorm:"column:description"`
	Table       string   `json:"table" gorm:"column:table"`
}

func NewRelation(data map[string]any, dbSchema *dbschema.DBSchema) *Relation {
	if dbSchema == nil {
		panic(errors.New("dbSch is nil"))
	}
	tenantId, _ := maputils.GetString(data, "tenant_id", "")
	caseId, _ := maputils.GetString(data, "case_id", "")
	id, _ := maputils.GetString(data, "id", "")
	relStartId, _ := getRelStartId(dbSchema, data)
	relEndId, _ := getRelEndId(dbSchema, data)
	relType, _ := getRelType(dbSchema, data)
	desc := utils.GetDescription(data, dbSchema)
	return &Relation{
		Id:          id,
		CaseId:      caseId,
		TenantId:    tenantId,
		Keywords:    []string{relType},
		RelType:     relType,
		Source:      relStartId,
		Target:      relEndId,
		SourceIds:   id,
		Table:       dbSchema.TableName,
		SourceType:  SourceType,
		Description: desc,
	}
}

func getRelStartId(dbSch *dbschema.DBSchema, data map[string]any) (string, error) {
	field := dbSch.GetRelStartIdField()
	if field == nil {
		return "", errors.New("start id not found")
	}
	val, err := maputils.GetString(data, field.DBName, "")
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
	val, err := maputils.GetString(data, field.DBName, "")
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
	val, err := maputils.GetString(data, field.DBName, "")
	if err != nil {
		return "", err
	}
	return val, nil
}
