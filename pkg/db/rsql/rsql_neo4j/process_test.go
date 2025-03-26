package rsql_neo4j

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ParseProcess(t *testing.T) {
	input := "((toto==32 and userId=='001' ) or (user=='admin' and sex==1)) and user=~'*000*' and userId=='001' "
	p := NewProcess("test", "n")
	err := rsql.ParseProcess(input, p)
	assert.NoError(t, err)
	sql := p.GetSQL()
	t.Log(sql)
}

func Test_InSubTable(t *testing.T) {

	t.Run("sub", func(t *testing.T) {
		input := "id=='1111'"
		p := NewProcess("test", "n")
		err := rsql.ParseProcess(input, p)
		assert.NoError(t, err)

		sql := p.GetSQL()
		t.Log(sql)
	})

	t.Run("sub", func(t *testing.T) {
		input := "id=in=sub(table:orderItems,field:customerId,rsql:product~='*book*')"
		p := NewProcess("test", "n")
		err := rsql.ParseProcess(input, p)
		assert.NoError(t, err)

		sql := p.GetSQL()
		t.Log(sql)
	})

	t.Run("in []string", func(t *testing.T) {
		input := "id=in=('a1','b2','c3')"
		p := NewProcess("test", "n")
		err := rsql.ParseProcess(input, p)
		assert.NoError(t, err)

		sql := p.GetSQL()
		t.Log(sql)
	})

	t.Run("in []number", func(t *testing.T) {
		input := "id=in=(12,22,31)"
		p := NewProcess("test", "n")
		err := rsql.ParseProcess(input, p)
		assert.NoError(t, err)

		sql := p.GetSQL()
		t.Log(sql)
	})
}
