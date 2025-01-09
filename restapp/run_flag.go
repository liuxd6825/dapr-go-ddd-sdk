package restapp

type RunType int

type RunFlag struct {
	RunType  RunType
	Env      string
	Config   string
	MainFile string // hServer主文件
	WorkPath string // hServer项目目录
	SqlFile  string
	DbKey    string
	Prefix   string
}

const (
	RunTypeStart RunType = iota
	RunTypeStop
	RunTypeStatus
	RunTypeInitDB
	RunTypeCreateSqlFile
	RunTypeVersion
	RunTypeHelp
)

var runFlags = &RunFlag{}

func GetRunFlag() *RunFlag {
	return runFlags
}
