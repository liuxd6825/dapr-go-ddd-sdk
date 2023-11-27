package restapp

import (
	"flag"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
)

type RunFlag struct {
	Help     bool
	Env      string
	Config   string
	InitDb   bool
	DbScript bool
	File     string
	DbKey    string
	Prefix   string
	Level    *logs.Level
}

func NewRunFlag(config string) *RunFlag {
	help := flag.Bool("help", false, "参数说明")
	env := flag.String("env", "", "环境变量名，可替换配置文件中的env值")
	cfg := flag.String("config", config, "配置文件")
	initdb := flag.Bool("initdb", false, "启动初始化模式，不启动服务")
	prefix := flag.String("prefix", "", "表前缀")
	lvl := flag.String("level", "", "日志级别")
	dbScript := flag.Bool("dbscript", false, "生成数据库脚本")
	file := flag.String("file", "", "脚本文件名称")
	dbKey := flag.String("dbKey", "", "")

	flag.Parse()

	flag := &RunFlag{
		Help:     *help,
		Env:      *env,
		Config:   *cfg,
		InitDb:   *initdb,
		Prefix:   *prefix,
		File:     *file,
		DbScript: *dbScript,
		DbKey:    *dbKey,
	}

	if *lvl != "" {
		if level, err := logs.ParseLevel(*lvl); err != nil {
			flag.Level = &level
		} else if err != nil {
			panic(err)
		}
	}

	return flag
}
