package engine

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSheetTemplate_Execute(t *testing.T) {
	ctx := context.Background()
	tpl := NewSheetTemplate(nil, nil, nil)
	sch := xtest.GetHumanSchema()
	html, err := tpl.Execute(ctx, sch)
	assert.Nil(t, err)
	assert.NotEmpty(t, html)
	t.Log(html)
}
