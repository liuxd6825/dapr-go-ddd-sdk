package envs

import (
	"strings"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/boolutils"
)

type GraphDBType string

const (
	GraphDBTypeNeo4j  DBType = "neo4j"
	GraphDBTypeHuge   DBType = "hugegraph"
	GraphDBTypeNebula DBType = "nebula"
)

type GraphEnv struct {
	Neo4j  map[string]*Neo4j
	Huge   map[string]*Huge
	Nebula map[string]*Nebula
}

func newGraphEnv() *GraphEnv {
	return &GraphEnv{
		Neo4j:  make(map[string]*Neo4j),
		Huge:   make(map[string]*Huge),
		Nebula: make(map[string]*Nebula),
	}
}

func initGraphDB(env *Env) {
	errList := errors.NewErrors()
	for dbKey, info := range env.Graphdb {
		dbType, isHas := info["dbType"]
		info["dbKey"] = strings.ToLower(dbKey)
		if !isHas {
			errList.AddError(errors.New("%s dbType is required", dbKey))
		}
		enable := boolutils.GetBool(info["enable"], true)
		if !enable {
			continue
		}

		switch dbType {
		case GraphDBTypeNeo4j.String():
			db, err := NewNeo4jWithMap(info)
			if err != nil {
				errList.AddError(errors.New("dbKey %s error: %s", dbKey, err.Error()))
				continue
			}
			if err := AddGraphDBNeo4j(env, db); err != nil {
				errList.AddError(err)
			}
			env.graphdb.Neo4j[db.DbKey] = db
			break
		case GraphDBTypeHuge.String():
			db, err := NewHugeWithMap(info)
			if err != nil {
				errList.AddError(errors.New("dbKey %s error: %s", dbKey, err.Error()))
				continue
			}
			if err := AddGraphDBHuge(env, db); err != nil {
				errList.AddError(err)
			}
			env.graphdb.Huge[db.DbKey] = db
			break
		case GraphDBTypeNebula.String():
			db, err := NewNebulaWithMap(info)
			if err != nil {
				errList.AddError(errors.New("dbKey %s error: %s", dbKey, err.Error()))
				continue
			}
			if err := AddGraphDBNebula(env, db); err != nil {
				errList.AddError(err)
			}
			env.graphdb.Nebula[db.DbKey] = db
			break
		default:
			errList.AddError(errors.New("dbType is invalid, ", dbKey))
		}
	}
	if errList.HasError() {
		panicError(errList)
	}
}
