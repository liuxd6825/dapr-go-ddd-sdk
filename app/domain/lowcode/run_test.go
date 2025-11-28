package lowcode

import (
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/lowcode/packages"
	"github.com/open2b/scriggo"
	"github.com/spf13/afero"
)

func Test_build(t *testing.T) {
	fs := afero.NewOsFs()

	packages := packages.GetPackages()
	program, err := build(fs, "./script", &scriggo.BuildOptions{
		Packages: packages,
	})
	
	if err != nil {
		t.Fatal(err)
	}
	err = program.Run(nil)
	if err != nil {
		t.Fatal(err)
	}
}
