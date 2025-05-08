package template

import (
	"bytes"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormBuilder_Build(t *testing.T) {
	s, err := schema.NewSchemaFile("./testfile/schema.json")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	ui, err := schema.NewUiSchemaFile(s, "./testfile/uischema.json")
	assert.NoError(t, err)
	assert.NotNil(t, ui)

	builder := NewUiBuilder(s, ui)
	var b bytes.Buffer
	err = builder.Build("./tpl/form.html", &b)
	assert.NoError(t, err)
	t.Log(b.String())
}
