package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/action"
)

type SuTaskService struct {
	tranDao    *dao.TranDao
	taskDao    *dao.SuTaskDao
	accountDao *dao.SuTaskAccountDao
}

func NewSuTaskService() *SuTaskService {
	return &SuTaskService{}
}

func (s *SuTaskService) Analyse(ctx context.Context, taskId string, rule model.SuTaskRule) error {
	task, err := s.taskDao.FindById(ctx, taskId)
	if err != nil {
		return err
	}
	accounts, err := s.accountDao.FindByTaskId(ctx, taskId)
	if err != nil {
		return err
	}
	analyse := action.NewSuTaskAnalyse(task, &task.SuTaskRule, accounts)
	analyse.Analyze2(ctx)
	return nil
}
