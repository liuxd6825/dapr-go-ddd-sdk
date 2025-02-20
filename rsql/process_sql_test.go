package rsql

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ParseProcess(t *testing.T) {
	input := "((toto==32 and userId=='001' ) or (user=='admin' and sex==1)) and user==~'000'"
	p := NewSqlProcess()
	err := ParseProcess(input, p)
	p.GetSQL()
	assert.Error(t, err)
}

func Test_InSubTable(t *testing.T) {
	input := "id=in=sub{table:orderItems;field:customerId;rsql:product~='*book*'}"
	//input := "id=in=sub(orderItems,customerId,product=like=*book*)"
	//.input := "id=in='001'"
	p := NewSqlProcess()
	err := ParseProcess(input, p)
	assert.NoError(t, err)

	sql := p.GetSQL()
	t.Log(sql)
}
