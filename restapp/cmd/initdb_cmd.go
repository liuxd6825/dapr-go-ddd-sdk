package cmd

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/spf13/cobra"
)

func newInitDbCmd() *cobra.Command {
	var initCmd = &cobra.Command{
		Use:   "init-db",
		Short: "初始化数据",
		Long:  "初始化数据库，建表、建字段、建索引等。",
		Run: func(cmd *cobra.Command, args []string) {
			runFlag.RunType = restapp.RunTypeInitDB
			if err := runFunc(runFlag); err != nil {
				panic(err)
			}
		},
	}
	initCmd.LocalFlags().StringVar(&runFlag.DbKey, "db-key", "", "")
	initCmd.LocalFlags().StringVar(&runFlag.Prefix, "prefix", "", "")
	return initCmd
}
