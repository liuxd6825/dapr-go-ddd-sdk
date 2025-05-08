package cmd

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/spf13/cobra"
	"os"
)

func newCreateSqlFileCmd() *cobra.Command {
	var initCmd = &cobra.Command{
		Use:   "sql-file",
		Short: "创建SQLFile",
		Long:  "生成数据库初始化脚本文件，。",
		Run: func(cmd *cobra.Command, args []string) {
			runFlag.RunType = restapp.RunTypeCreateSqlFile
			if err := runFunc(runFlag); err != nil {
				fmt.Println(err.Error())
				os.Exit(0)
			}
		},
	}
	initCmd.LocalFlags().StringVarP(&runFlag.SqlFile, "file", "f", "./init_db.sql", "要生成的初始化脚本文件名")
	return initCmd
}
