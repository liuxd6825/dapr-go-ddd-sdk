package restapp

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func Restart() error {
	// 获取当前可执行文件路径
	executable, err := os.Executable()
	if err != nil {
		fmt.Println("获取可执行文件路径失败:", err)
		return err
	}

	// 获取命令行参数
	args := os.Args[1:]

	// 构造重启命令
	cmd := exec.Command(executable, args...)
	cmd.Dir = "./" // 如果需要设置特定的启动目录，可以在此处修改启动目录

	// 设置重启时的环境变量，类似于当前进程的环境变量
	cmd.Env = os.Environ()

	// 启动新进程
	err = cmd.Start()
	if err != nil {
		fmt.Println("重启失败:", err)
		return err
	}

	// 退出当前进程，确保新的进程启动
	fmt.Println("正在重启...")
	// 通过 syscall.Exit(0) 确保退出当前进程
	syscall.Exit(0)
	return nil
}
