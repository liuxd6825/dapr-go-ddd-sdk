package rsql_mongo

import (
	"encoding/json"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rsql"
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

func Test_SubQuery(t *testing.T) {
	process(t, "test", "id=in=sub(table:human_certificate, field:human_id, rsql:certificate_code~='2222')")
}

func process(t *testing.T, tenantId string, input string) any {
	p := NewProcess(tenantId)
	err := rsql.ParseProcess(input, p)
	assert.Error(t, err)
	data := p.GetFilter()

	jsonText, err := json.Marshal(data)
	assert.NoError(t, err)
	fmt.Println(string(jsonText))
	return data
}
