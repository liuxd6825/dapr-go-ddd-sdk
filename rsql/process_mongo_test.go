package rsql

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_mongodb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Process1(t *testing.T) {
	process(t, "001", "name=='A' ")
}

func TestProcess2(t *testing.T) {
	process(t, "001", "name=='A' and name=='B'")
}

func TestProcess3(t *testing.T) {
	process(t, "001", "(name=='A' and name=='B') or (name=='C')")
}

func TestIsNull(t *testing.T) {
	process(t, "001", "taskId=='0001' and errors=!null=0")
}

func process(t *testing.T, tenantId string, input string) map[string]interface{} {
	p := ddd_mongodb.NewMongoProcess()
	err := ParseProcess(input, p)
	assert.Error(t, err)
	data, err := p.GetFilter(tenantId)
	return data
}
