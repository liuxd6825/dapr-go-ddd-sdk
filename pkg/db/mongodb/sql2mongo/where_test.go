package sql2mongo

import (
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
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
