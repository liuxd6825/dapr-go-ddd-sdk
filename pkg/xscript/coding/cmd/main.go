package main

import (
	"flag"
	"os"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xscript/coding/packages"
)

func main() {
	// 命令行参数
	srcDir := flag.String("src", "", "根源目录路径 (必须)")
	path := flag.String("path", "", "根目录对应的基础导入路径 (必须)")
	outputFile := flag.String("out", "scriggo_exports.go", "输出文件名")
	outPkgName := flag.String("pkg", "exports", "生成文件的包名")
	recursive := flag.Bool("recursive", false, "是否递归扫描子目录")
	flag.Parse()

	if *srcDir == "" || *path == "" {
		flag.Usage()
		os.Exit(1)
	}
	err := packages.Make(*srcDir, *path, *outPkgName, *outputFile, *recursive)
	if err != nil {
		panic(err)
	}
}
