package restapp

import (
	"flag"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
)

type RunFlag struct {
	Help    bool
	Env     string
	Config  string
	Init    bool
	SqlFile string
	DbKey   string
	Prefix  string
	Level   *logs.Level
}

func NewRunFlag(config string) *RunFlag {
	help := flag.Bool("help", false, "参数说明")
	env := flag.String("env", "", "环境变量名，可替换配置文件中的env值")
	cfg := flag.String("config", config, "配置文件")
	init := flag.Bool("init", false, "初始化数据库")
	prefix := flag.String("prefix", "", "表前缀")
	lvl := flag.String("level", "", "日志级别")
	sqlFile := flag.String("sqlfile", "", "生成数据库脚本的文件名")
	dbKey := flag.String("dbkey", "", "")

	flag.Parse()

	flag := &RunFlag{
		Help:    *help,
		Env:     *env,
		Config:  *cfg,
		Init:    *init,
		Prefix:  *prefix,
		SqlFile: *sqlFile,
		DbKey:   *dbKey,
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
