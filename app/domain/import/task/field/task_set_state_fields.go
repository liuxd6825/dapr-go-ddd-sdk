package field

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/task/enums"
)

type TaskUpdateStateFields struct {
	Id      string          `json:"id" vali`
	State   enums.TaskState `json:"state"`
	Message string          `json:"message"`
}
