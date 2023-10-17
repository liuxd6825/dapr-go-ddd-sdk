package restapp

import (
	"flag"
)

type RunFlag struct {
	Help   bool
	Env    string
	Config string
	Init   bool
	Prefix string
}

func NewRunFlag(config string) *RunFlag {
	help := flag.Bool("help", false, "参数说明")
	env := flag.String("env", "", "环境变量名，可替换配置文件中的env值")
	cfg := flag.String("config", config, "配置文件")
	init := flag.Bool("init", false, "启动初始化模式，不启动服务")
	prefix := flag.String("prefix", "", "表前缀")

	flag.Parse()

	return &RunFlag{
		Help:   *help,
		Env:    *env,
		Config: *cfg,
		Init:   *init,
		Prefix: *prefix,
	}
}
