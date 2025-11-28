package lowcode

import (
	testing2 "testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xscript/coding/packages"
)

func Test_MakeCode(t *testing2.T) {
	if err := packages.Make(packages.MakeOption{
		SrcDir:     "../../pkg",
		ImportPath: "github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg",
		OutPkgName: "pkg",
		OutputFile: "./packages/app/pkg/pkg.go",
		Recursive:  true,
	}); err != nil {
		panic(err)
	}

	if err := packages.Make(packages.MakeOption{
		SrcDir:     "../../domain",
		ImportPath: "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain",
		OutPkgName: "domain",
		OutputFile: "./packages/app/domain/domain.go",
		Recursive:  true,
		Cancels: []string{
			"app/domain/lowcode",
		},
	}); err != nil {
		panic(err)
	}

	if err := packages.Make(packages.MakeOption{
		SrcDir:     "../../../pkg",
		ImportPath: "github.com/liuxd6825/dapr-go-ddd-sdk/pkg",
		OutPkgName: "pkg",
		OutputFile: "./packages/pkg/pkg.go",
		Recursive:  true,
	}); err != nil {
		panic(err)
	}

}
