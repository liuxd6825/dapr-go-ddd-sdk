package service

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/command"
	dao2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/field"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	doc "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/outside"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/service/interfaces"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/singleutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/timeutils"
)

type TaskService struct {
	dao *dao2.TaskDao
	xbase.Service
	docStatusProvider interfaces.ImportStatusProvider
}

func NewTaskService() *TaskService {
	return singleutils.CreateObj[*TaskService](func() *TaskService {
		taskService := &TaskService{
			dao:               dao2.NewTaskDao(config.DBKey),
			docStatusProvider: outside.NewImportStatusProvider(),
		}
		return taskService
	})
}

func (r *TaskService) Create(ctx context.Context, cmd *command.TaskCreateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		entity := &model.Task{}
		entity.Id = cmd.Data.Id
		entity.CaseId = cmd.Data.CaseId
		entity.DocId = cmd.Data.DocId
		entity.State = model.TaskStateEditing
		entity.SheetName = cmd.Data.SheetName
		entity.SheetId = cmd.Data.SheetId
		entity.FileName = cmd.Data.FileName
		entity.FileId = cmd.Data.FileId
		entity.Remark = cmd.Data.Remark
		entity.SchemaId = cmd.Data.SchemaId
		entity.SchemaName = cmd.Data.SchemaName
		entity.MapHeads = cmd.Data.MapHeads
		entity.Fields = cmd.Data.Fields
		entity.MasterType = cmd.Data.MasterType
		entity.MasterId = cmd.Data.MasterId
		entity.MasterName = cmd.Data.MasterName
		err := r.dao.Create(ctx, entity).GetError()
		if err != nil {
			return err
		}

		docModel := &doc.Document{}
		docModel.Id = idutils.NewId()
		docModel.FileId = entity.FileId
		docModel.SourceId = entity.DocId
		docModel.SourceType = "流水"
		docModel.CaseId = entity.CaseId
		docModel.SourceApp = "document_service"
		err = r.docStatusProvider.UpdateStatus(ctx, docModel, model.TaskStateEditing.Name(), "")
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *TaskService) Update(ctx context.Context, cmd *command.TaskUpdateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		entity := &model.Task{}
		entity.Id = cmd.Data.Id
		entity.DocId = cmd.Data.DocId
		entity.SheetId = cmd.Data.SheetId
		entity.SheetName = cmd.Data.SheetName
		entity.FileName = cmd.Data.FileName
		entity.FileId = cmd.Data.FileId
		entity.Remark = cmd.Data.Remark
		entity.MapHeads = cmd.Data.MapHeads
		entity.Fields = cmd.Data.Fields
		entity.MasterType = cmd.Data.MasterType
		entity.MasterId = cmd.Data.MasterId
		entity.MasterName = cmd.Data.MasterName

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
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
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
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return r.dao.UpdateLock(ctx, cmd.Data.Id, false)
	})
}

func (r *TaskService) UpdateStart(ctx context.Context, cmd *command.TaskRecordCreateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		data := map[string]any{
			"start_time": timeutils.PNow(),
			//"state":      model.TaskStateStart,
		}
		return r.dao.UpdateMap(ctx, cmd.Data.Id, data).GetError()
	})
}

func (r *TaskService) Stop(ctx context.Context, cmd *command.TaskStopCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		data := map[string]any{
			"end_time": timeutils.PNow(),
			//"state":    model.TaskStateStop,
		}
		return r.dao.UpdateMap(ctx, cmd.Data.Id, data).GetError()
	})
}

func (r *TaskService) Delete(ctx context.Context, cmd *command.TaskDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		id := cmd.Data.Id
		task, err := r.dao.FindById(ctx, id)
		if err != nil {
			return err
		}

		switch task.State {
		case model.TaskStateEditing, model.TaskStateGenerated, model.TaskStateNone:
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
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return r.SetState(ctx, cmd.Data.Id, cmd.Data.State, cmd.Data.Message)
	})
}

func (r *TaskService) SetState(ctx context.Context, taskId string, state model.TaskState, message string) error {
	return r.dao.SetState(ctx, taskId, state, message)
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
	qry.SetMustFilter(fmt.Sprintf("case_id=='%s'", qry.GetCaseId()))
	res := r.dao.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (r *TaskService) FindByFileId(ctx context.Context, caseId string, fileId string) (*model.Task, error) {
	task, err := r.dao.FindOneByRSQL(ctx, fmt.Sprintf("case_id=='%s' and file_id=='%s'", caseId, fileId))
	return task, err
}
