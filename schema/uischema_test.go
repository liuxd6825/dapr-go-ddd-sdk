package schema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUiSchemaFile(t *testing.T) {
	uiSchema, err := NewUiSchemaFile("./test/uischema.json")
	assert.NoError(t, err)
	assert.NotNil(t, uiSchema)

	schema, err := NewSchemaFile("./test/schema.json")
	assert.NoError(t, err)
	assert.NotNil(t, schema)

}
