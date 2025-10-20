package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type LoginCommand struct {
	xbase.Command[model.Login]
}
