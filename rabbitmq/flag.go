package rabbitmq

import "flag"

type Flag struct {
	Help   bool
	Host   string
	User   string
	Pwd    string
	Reset  string
	Filter string
	Cmd    string
}

func NewFlag() *Flag {
	help := flag.Bool("help", false, "参数说明")
	host := flag.String("host", "127.0.0.1:15672", "环境变量名，可替换配置文件中的env值")
	user := flag.String("user", "rabbitmq", "用户名")
	pwd := flag.String("pwd", "rabbitmq_15672", "密码")
	filter := flag.String("filter", "", "名称过滤条件")
	cmd := flag.String("cmd", "", "命令类型，clearExchange|clearQueue")

	flag.Parse()

	f := &Flag{
		Help:   *help,
		Host:   *host,
		User:   *user,
		Pwd:    *pwd,
		Filter: *filter,
		Cmd:    *cmd,
	}
	return f
}
