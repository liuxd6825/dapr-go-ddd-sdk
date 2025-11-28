package xscript

import (
	"fmt"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/open2b/scriggo/native"
)

func AddPackages(target native.Packages, src native.Packages) {
	for key, pkg := range src {
		target[key] = pkg
	}
}

func AddDeclarations(target native.Declarations, src native.Declarations) {
	for key, pkg := range src {
		target[key] = pkg
	}
}

// GetStandardPackages
// @Description: 加载go包
// @param opts
// @return native.Packages
// @return error
func GetStandardPackages() native.Packages {
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
