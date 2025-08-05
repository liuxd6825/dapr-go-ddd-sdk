package query

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

// RecordIeFindPagingByTaskIdQuery 分页查询命令
type RecordIeFindPagingByTaskIdQuery struct {
	idao.FindPagingQueryRequest
	TaskId      string `json:"taskId" param:"task-id" `
	IsFindError bool   `json:"findError" param:"is-find-error"`
}

type RecordFindByIdQueryRequest = store.FindByIdQueryRequest

func NewRecordIeFindPagingByTaskIdQuery(taskId string, findError bool) *RecordIeFindPagingByTaskIdQuery {
	qry := &RecordIeFindPagingByTaskIdQuery{
		TaskId:      taskId,
		IsFindError: findError,
	}
	qry.TaskId = taskId
	qry.IsTotalRows = true
	qry.PageSize = 1000
	qry.PageNum = 1
	qry.Filter = ""
	return qry
}
