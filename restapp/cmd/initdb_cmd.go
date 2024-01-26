package cmd

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/spf13/cobra"
	"os"
)

func newInitDbCmd() *cobra.Command {
	var initCmd = &cobra.Command{
		Use:   "init-db",
		Short: "初始化数据",
		Long:  "初始化数据库，建表、建字段、建索引等。",
		Args:  cobra.MatchAll(cobra.ExactArgs(0)),
		Run: func(cmd *cobra.Command, args []string) {
			runFlag.RunType = restapp.RunTypeInitDB
			if err := runFunc(runFlag); err != nil {
				fmt.Println(err.Error())
				os.Exit(0)
			}
		},
	}
	initCmd.LocalFlags().StringVarP(&runFlag.DbKey, "db-key", "d", "", "配置文件中数据库的关键字")
	initCmd.LocalFlags().StringVarP(&runFlag.Prefix, "prefix", "o", "", "数据表名的前缀符")
	return initCmd
}
