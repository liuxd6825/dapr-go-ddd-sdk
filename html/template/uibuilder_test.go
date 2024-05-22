package template

import (
	"bytes"
	"github.com/liuxd6825/dapr-go-ddd-sdk/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormBuilder_Build(t *testing.T) {
	s, err := schema.NewSchemaFile("./test/schema.json")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	ui, err := schema.NewUiSchemaFile(s, "./test/uischema.json")
	assert.NoError(t, err)
	assert.NotNil(t, ui)

	builder := NewUiBuilder(s, ui)
	var b bytes.Buffer
	err = builder.Build("./tpl/form.html", &b)
	assert.NoError(t, err)
	t.Log(b.String())
}
