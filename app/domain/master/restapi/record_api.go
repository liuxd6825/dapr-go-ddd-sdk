package restapi

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service"

type RecordAPI struct {
	service *service.RecordService
}

func NewRecordAPI() *RecordAPI {
	return &RecordAPI{
		service: service.NewRecordService(),
	}
}
