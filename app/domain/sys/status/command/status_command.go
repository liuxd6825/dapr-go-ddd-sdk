package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/status/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type StatusCreateCommand = xbase.Command[*model.Status]

type StatusUpdateCommand = xbase.Command[*model.Status]
