package render

import (
	"bytes"
	"fmt"
	"os/exec"
)

// CompileTS 通过esbuild直接编译TS代码，返回JS代码
func CompileTS(tsCode []byte) ([]byte, error) {
	cmd := exec.Command("esbuild", "--loader=ts", "--format=cjs")

	// 创建输入缓冲区
	var stdin bytes.Buffer
	stdin.Write(tsCode)
	cmd.Stdin = &stdin

	// 捕获输出
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 执行命令
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("编译失败: %v\n错误输出: %s", err, stderr.String())
	}
	return stdout.Bytes(), nil
}
