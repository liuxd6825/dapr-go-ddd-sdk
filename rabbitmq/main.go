package rabbitmq

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
)

func Main() {
	flag := NewFlag()
	if flag.Help {
		return
	}
	m := NewManage(flag.Host, flag.User, flag.Pwd)
	var err error
	switch flag.Cmd {
	case "clearExchange":
		{
			err = m.ClearAllExchange(flag.Filter)
		}
	case "clearQueue":
		{
			err = m.ClearAllQueues(flag.Filter)
		}
	case "clear":
		{
			err = m.ClearAllExchange(flag.Filter)
			if err != nil {
				break
			}
			err = m.ClearAllQueues(flag.Filter)
		}
	default:
		err = errors.New("error : cmd is clearExchange clearQueue clear ")
	}
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("OK")
}
