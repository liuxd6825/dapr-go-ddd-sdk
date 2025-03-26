package entity

import (
	"github.com/liuxd6825/jsonschema/v6"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Schema(t *testing.T) {
	sch, err := jsonschema.NewSchemaWithStruct("schema.json", &User{})
	assert.NoError(t, err)
	t.Log("schema:", sch)
}
