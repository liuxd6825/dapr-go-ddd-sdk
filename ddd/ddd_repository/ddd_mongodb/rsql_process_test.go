package ddd_mongodb

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/rsql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Process1(t *testing.T) {
	if _, err := process(t, "001", "name=='A' "); err != nil {
		t.Error(err)
	}
}

func TestProcess2(t *testing.T) {
	if _, err := process(t, "001", "name=='A' and name=='B'"); err != nil {
		t.Error(err)
	}
}

func TestProcess3(t *testing.T) {
	if _, err := process(t, "001", "(name=='A' and name=='B') or (name=='C')"); err != nil {
		t.Error(err)
	}
}

func process(t *testing.T, tenantId string, input string) (map[string]interface{}, error) {
	p := NewMongoProcess()
	err := rsql.ParseProcess(input, p)
	assert.Error(t, err)
	data, err := p.GetFilter(tenantId)
	if err != nil {
		return nil, err
	}
	return data, nil
}
