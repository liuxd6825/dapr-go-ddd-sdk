package sql2mongo

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"testing"
)

func Test_Parse(t *testing.T) {
	gp.Try(func() error {
		opts := Parse("select count(*) from users")
		t.Log(" filter: ", opts.Filter, " opts: ", opts)
		return nil
	}).Catch(func(e error) {
		t.Error(e)
	})

}
