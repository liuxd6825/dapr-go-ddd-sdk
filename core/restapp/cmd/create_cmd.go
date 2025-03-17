package cmd

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/spf13/cobra"
	"os"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "创建",
	Long:  "创建HTML文件",
	Args:  cobra.MatchAll(cobra.ExactArgs(0)),
	Run: func(cmd *cobra.Command, args []string) {
		runFlag.RunType = restapp.RunTypeStart
		if err := runFunc(runFlag); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	},
}

var createType string
var createSchema string
var createFile string

func getCreateCmd() *cobra.Command {
	createCmd.PersistentFlags().StringVar(&createType, "tpl", "form", "配置文件名")
	createCmd.PersistentFlags().StringVar(&createSchema, "schema", "", "Schema配置文件")
	createCmd.PersistentFlags().StringVar(&createFile, "file", "", "要生成的文件名称")
	return startCmd
}
