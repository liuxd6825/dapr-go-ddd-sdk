package service

import (
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
)

type AccountService struct {
	dao *dao.RecordDao
}

var _accountService *AccountService
var _accountServiceOnce sync.Once

func NewAccountService() *AccountService {
	_accountServiceOnce.Do(func() {
		_accountService = newAccountService()
	})
	return _accountService
}

func newAccountService() *AccountService {
	return &AccountService{
		dao: dao.NewRecordDao(config.DBKey),
	}
}
