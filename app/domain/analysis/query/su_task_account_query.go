package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"

// SuTaskAccountFindByIdQuery 按聚合根ID查询命令
type SuTaskAccountFindByIdQuery struct {
	Id string `json:"id" param:"id"`
}

type TaskAccountFindPagingByTaskIdQuery struct {
	TaskId string `json:"taskId"  param:"taskId" required:"true"`
	idao.FindPagingQueryRequest
}
