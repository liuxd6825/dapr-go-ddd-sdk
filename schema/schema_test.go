package schema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSchemaString(t *testing.T) {
	schema, err := NewSchemaFile("./test/schema.json")
	assert.NoError(t, err)
	assert.NotNil(t, schema)
}
