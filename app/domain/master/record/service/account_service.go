package service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/dao"
)

type AccountService struct {
	dao *dao.RecordDao
}

func NewAccountService(dbKey string) *AccountService {
	service := newAccountService(dbKey)
	return service
}

func newAccountService(dbKey string) *AccountService {
	return &AccountService{
		dao: dao.NewRecordDao(dbKey),
	}
}
