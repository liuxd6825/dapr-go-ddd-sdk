package schema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUiSchemaFile(t *testing.T) {
	schema, err := NewSchemaFile("./test/validate.json")
	assert.NoError(t, err)
	assert.NotNil(t, schema)

	uiSchema, err := NewUiSchemaFile(schema, "./test/uischema.json")
	assert.NoError(t, err)
	assert.NotNil(t, uiSchema)

}
