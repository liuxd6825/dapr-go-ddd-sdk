package field

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/task/enums"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
)

type TaskUpdateProgressFields struct {
	Id           string          `json:"id" title:"开始时间"`
	CompleteRows int64           `json:"completeRows" title:"开始时间"`
	State        enums.TaskState `json:"state"  title:"状态"`
	Message      string          `json:"message"  title:"开始时间"`
	StartTime    *times.Time     `json:"startTime"   title:"开始时间"`
	EndTime      *times.Time     `json:"endTime"  title:"最后日期"`
}
