package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/recordie/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/task/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/task/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/task/enums"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/task/field"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/task/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/singleutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/timeutils"
)

type TaskService struct {
	dao *dao.TaskDao
	xbase.Service
}

func NewTaskService() *TaskService {
	return singleutils.CreateObj[*TaskService](func() *TaskService {
		return &TaskService{
			dao: dao.NewTaskDao(config.DBKey),
		}
	})
}

func (r *TaskService) Create(ctx context.Context, cmd *command.TaskCreateCommand) error {
	return r.DoCommand(ctx, cmd, func(ctx context.Context) error {
		entity := &model.Task{}
		entity.Id = cmd.Data.Id
		entity.TenantId = cmd.Data.TenantId
		entity.CaseId = cmd.Data.CaseId
		entity.DocId = cmd.Data.DocId
		entity.State = enums.TaskStateEditing
		entity.SheetName = cmd.Data.SheetName
		entity.FileName = cmd.Data.FileName
		entity.FileId = cmd.Data.FileId
		entity.Remark = cmd.Data.Remark
		entity.SchemaId = cmd.Data.SchemaId
		entity.MapHeads = cmd.Data.MapHeads
		entity.Fields = cmd.Data.Fields
		return r.dao.Create(ctx, entity).GetError()
	})

}

func (r *TaskService) Update(ctx context.Context, cmd *command.TaskUpdateCommand) error {
	return r.DoCommand(ctx, cmd, func(ctx context.Context) error {
		entity := &model.Task{}
		entity.Id = cmd.Data.Id
		entity.TenantId = cmd.Data.TenantId
		entity.DocId = cmd.Data.DocId
		entity.SheetName = cmd.Data.SheetName
		entity.FileName = cmd.Data.FileName
		entity.FileId = cmd.Data.FileId
		entity.Remark = cmd.Data.Remarks
		entity.MapHeads = cmd.Data.MapHeads
		entity.Fields = cmd.Data.Fields

		return r.dao.Update(ctx, entity, idao.NewCallOptions().SetUpdateFields(cmd.UpdateMask)).GetError()
	})
}

func (r *TaskService) CheckLock(ctx context.Context, taskId string) (bool, error) {
	task, err := r.FindById(ctx, taskId)
	if err != nil {
		return false, err
	} else if task == nil {
		return false, errors.ErrDbNotFoundRecord(taskId)
	}
	if task.Lock {
		return false, errors.ErrDbRecordLock()
	}
	return true, nil
}

func (r *TaskService) Lock(ctx context.Context, cmd *command.TaskLockCommand) error {
	return r.DoCommand(ctx, cmd, func(ctx context.Context) error {
		task, err := r.FindById(ctx, cmd.Data.Id)
		if err != nil {
			return err
		} else if task == nil {
			return errors.ErrDbNotFoundRecord(cmd.Data.Id)
		}
		if task.Lock {
			return errors.ErrDbRecordLock()
		}
		return r.dao.UpdateLock(ctx, cmd.Data.Id, true)
	})
}

func (r *TaskService) Unlock(ctx context.Context, cmd *command.TaskUnlockCommand) error {
	return r.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return r.dao.UpdateLock(ctx, cmd.Data.Id, false)
	})
}

func (r *TaskService) UpdateStart(ctx context.Context, cmd *command.TaskRecordCreateCommand) error {
	return r.DoCommand(ctx, cmd, func(ctx context.Context) error {
		data := map[string]any{
			"start_time": timeutils.PNow(),
			//"state":      model.TaskStateStart,
		}
		return r.dao.UpdateMap(ctx, cmd.Data.Id, data).GetError()
	})
}

func (r *TaskService) Stop(ctx context.Context, cmd *command.TaskStopCommand) error {
	return r.DoCommand(ctx, cmd, func(ctx context.Context) error {
		data := map[string]any{
			"end_time": timeutils.PNow(),
			//"state":    model.TaskStateStop,
		}
		return r.dao.UpdateMap(ctx, cmd.Data.Id, data).GetError()
	})
}

func (r *TaskService) Validate(ctx context.Context, m *command.TaskValidateCommand) error {
	return nil
}

func (r *TaskService) Delete(ctx context.Context, cmd *command.TaskDeleteCommand) error {
	return r.DoCommand(ctx, cmd, func(ctx context.Context) error {
		id := cmd.Data.Id
		task, err := r.dao.FindById(ctx, id)
		if err != nil {
			return err
		}

		switch task.State {
		case enums.TaskStateEditing, enums.TaskStateGenerated, enums.TaskStateNone:
			err = r.dao.DeleteById(ctx, id).GetError()
		default:
			err = errors.New("状态不正确%s", task.State.Name())
		}
		return err
	})
}

func (r *TaskService) UpdateProgress(ctx context.Context, progress *field.TaskUpdateProgressFields) error {
	return r.dao.UpdateProgress(ctx, progress)
}

func (r *TaskService) UpdateState(ctx context.Context, cmd *command.TaskUpdateStateCommand) error {
	return r.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return r.dao.SetState(ctx, cmd.Data.Id, cmd.Data.State, cmd.Data.Message)
	})
}

func (r *TaskService) FindById(ctx context.Context, id string) (*model.Task, error) {
	task, err := r.dao.FindById(ctx, id)
	return task, err
}

func (r *TaskService) FindPaging(ctx context.Context, qry idao.FindPagingQuery) (idao.FindPagingResult[*model.Task], error) {
	res := r.dao.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (r *TaskService) FindPagingByCaseId(ctx context.Context, qry idao.FindPagingByCaseIdQuery) (idao.FindPagingResult[*model.Task], error) {
	res := r.dao.FindPaging(ctx, qry)
	return res, res.GetError()
}
