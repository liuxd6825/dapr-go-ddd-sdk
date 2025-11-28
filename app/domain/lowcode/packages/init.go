package packages

import (
	"fmt"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/lowcode/packages/app/domain"

	/*domain "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/lowcode/packages/app/domain"
	app_pkg "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/lowcode/packages/app/pkg"
	pkg "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/lowcode/packages/pkg"*/
	"github.com/open2b/scriggo/native"
)

func GetPackages() native.Packages {
	packages := standard()
	addPackages(packages, domain.Packages())
	/*	addPackages(packages, app_pkg.Packages())
		addPackages(packages, pkg.Packages())*/
	return packages
}

func addPackages(target native.Packages, src native.Packages) {
	for key, pkg := range src {
		target[key] = pkg
	}
}

// LoadPackages
// @Description: 加载go包
// @param opts
// @return native.Packages
// @return error
func standard() native.Packages {
	// 1. 定义包映射
	packages := native.Packages{
		// 模拟 "fmt" 包
		"fmt": native.Package{
			Name: "fmt",
			Declarations: native.Declarations{
				"Println":  fmt.Println,
				"Printf":   fmt.Printf,
				"Sprintf":  fmt.Sprintf,
				"Sprint":   fmt.Sprint,
				"Sprintln": fmt.Sprintln,
				"Fprintf":  fmt.Fprintf,
				"Fprintln": fmt.Fprintln,
			},
		},
		"strings": native.Package{
			Name: "strings",
			Declarations: native.Declarations{
				"Split":   strings.Split,
				"ToUpper": strings.ToUpper,
			},
		},
		"iris": native.Package{
			Name: "iris",
			Declarations: native.Declarations{
				"Application": (*iris.Application)(nil),
				"Context":     (*iris.Context)(nil),
			},
		},
	}
	return packages
}
