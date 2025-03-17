package rsql

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/rsql/rsql_sql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_NewBuilder(t *testing.T) {

	t.Run("and", func(t *testing.T) {
		b := NewBuilder()
		// 简单AND查询
		query := b.And(
			b.Eq("name", "Alice"),
			b.Gt("age", 25),
		).Build()

		t.Log("rsql:", query) // (name=="Alice";age=gt=25)
		p := rsql_sql.NewProcess("test")
		err := ParseProcess(query, p)
		assert.NoError(t, err)
		t.Log("sql:", p.GetSQL())
	})

	t.Run("复杂嵌套查询", func(t *testing.T) {
		b := NewBuilder()
		// 复杂嵌套查询
		query := b.And(
			b.Eq("department", "engineering"),
			b.Or(
				b.And(
					b.Ge("salary", 50000),
					b.Le("salary", 100000),
				),
				b.In("role", []string{"manager", "lead"}),
			),
			b.Neq("status", "inactive"),
		).Build()
		// (department=="engineering";(salary=ge=50000;salary=le=100000,role=in=("manager","lead"));status!="inactive")
		t.Log("rsql:", query)

		p := rsql_sql.NewProcess("test")
		err := ParseProcess(query, p)
		assert.NoError(t, err)
		t.Log("sql:", p.GetSQL())
	})

	t.Run("Gt In", func(t *testing.T) {
		b := NewBuilder()
		// 处理特殊类型
		vt := time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC)
		query := b.And(
			b.Gt("createdAt", vt),
			b.In("tags", []string{"golang", "backend"}),
			b.Out("tags", []string{"java", "typescript"}),
		).Build()
		t.Log("rsql:", query)
		// (createdAt=gt="2023-10-05T00:00:00Z";tags=in=("golang","backend"))

		p := rsql_sql.NewProcess("test")
		err := ParseProcess(query, p)
		assert.NoError(t, err)
		t.Log("sql:", p.GetSQL())
	})

	t.Run("处理空值", func(t *testing.T) {
		b := NewBuilder()
		// 处理空值
		query := b.Or(
			b.Null("deletedAt"),
			b.NotNull("archived"),
		).Build()
		// (deletedAt==null,archived==null)

		p := rsql_sql.NewProcess("test")
		err := ParseProcess(query, p)
		assert.NoError(t, err)
		t.Log("rsql:", query)
		t.Log("sql:", p.GetSQL())
	})

}
