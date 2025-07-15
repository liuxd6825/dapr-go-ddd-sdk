package field

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/recordie/model"
)

type TaskUpdateRecordTemplateFields struct {
	Id       string               `json:"id"`
	TenantId string               `json:"tenantId"`
	Template model.RecordTemplate `json:"template"`
}
