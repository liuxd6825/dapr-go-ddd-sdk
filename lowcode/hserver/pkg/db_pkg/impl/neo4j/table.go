package neo4j

import (
	"github.com/liuxd6825/jsonschema/v6"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Table struct {
	tableName string
	schema    *jsonschema.Schema
	driver    *neo4j.DriverWithContext
}
