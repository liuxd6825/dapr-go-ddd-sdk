package runtime

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
)

type Compiler struct {
}

func NewCompiler() *Compiler {
	return &Compiler{}
}

func (c *Compiler) compileTypeScript(workingDir, tsCode string) (string, error) {
	// 创建动态的 transpile.js 脚本
	scriptPath, err := c.createTranspileScript(workingDir)
	if err != nil {
		return "", err
	}
	//defer os.Remove(scriptPath) // 确保脚本文件被删除

	// 调用 Node.js 脚本
	cmd := exec.Command("node", scriptPath)

	cmd.Dir = workingDir // 设置 Node.js 的工作目录
	// 将 TypeScript 代码写入标准输入
	cmd.Stdin = bytes.NewReader([]byte(tsCode))

	// 捕获输出
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// 执行命令
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to execute transpile.js: %v\nstderr: %s", err, stderr.String())
	}

	return out.String(), nil
}

func (c *Compiler) createTranspileScript(workingDir string) (string, error) {
	// 动态生成的 Node.js 脚本内容
	scriptContent := `
		const ts = require("typescript");

		process.stdin.on("data", (data) => {
			const tsCode = data.toString();
			try {
				// 使用 TypeScript API 转换代码
				const result = ts.transpileModule(tsCode, {
					compilerOptions: { module: ts.ModuleKind.CommonJS }
				});
				console.log(result.outputText);
			} catch (error) {
				console.error("Error during transpilation:", error.message);
				process.exit(1);
			}
		});
	`

	// 创建临时文件
	tempFile, err := ioutil.TempFile(workingDir, "transpile-*.js")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}

	// 写入脚本内容
	if _, err := tempFile.Write([]byte(scriptContent)); err != nil {
		return "", fmt.Errorf("failed to write to temp file: %v", err)
	}

	// 关闭文件
	if err := tempFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temp file: %v", err)
	}

	return tempFile.Name(), nil
}
