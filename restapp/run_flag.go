package restapp

import (
	"flag"
)

type RunType int

type RunFlag struct {
	RunType RunType
	Env     string
	Config  string
	SqlFile string
	DbKey   string
	Prefix  string
}

const (
	RunTypeStart RunType = iota
	RunTypeStop
	RunTypeStatus
	RunTypeInitDB
	RunTypeCreateSqlFile
	RunTypeHelp
)

func NewRunFlag(config string) *RunFlag {
	env := flag.String("env", "", "环境变量名，可替换配置文件中的env值")
	cfg := flag.String("config", config, "配置文件")
	prefix := flag.String("prefix", "", "表前缀")
	sqlFile := flag.String("sqlfile", "", "生成数据库脚本的文件名")
	dbKey := flag.String("dbkey", "", "")

	flag.Parse()

	flag := &RunFlag{
		Env:     *env,
		Config:  *cfg,
		Prefix:  *prefix,
		SqlFile: *sqlFile,
		DbKey:   *dbKey,
	}
	return flag
}
