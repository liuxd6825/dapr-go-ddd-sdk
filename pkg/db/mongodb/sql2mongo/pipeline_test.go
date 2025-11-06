package sql2mongo

import (
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
)

func Test_Pipeline(t *testing.T) {
	gp.Try(func() error {
		opts, err := Pipeline("select count(*),sum(money) from users where name like 'lxd%'")
		t.Log(" pipeline: ", opts)
		return err
	}).Catch(func(e error) {
		t.Error(e)
	})

}
