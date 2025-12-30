package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dapr/go-sdk/client"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/outside"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/service/interfaces"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/doc_extract"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/entity"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/storage"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
)

type DocumentService struct {
	dao            *dao.DocumentDao
	graphRag       *my_rag.GraphRag
	docExtract     *doc_extract.Extract
	statusProvider interfaces.ImportStatusProvider
	current        *model.Document
}

var documentService *DocumentService
var documentServiceOnce sync.Once

func NewDocumentService(statusProvider interfaces.ImportStatusProvider) *DocumentService {
	documentServiceOnce.Do(func() {
		documentService = &DocumentService{
			dao:            dao.NewDocumentDao(DBKey),
			graphRag:       NewGraphRag(),
			docExtract:     doc_extract.NewExtract(logs.GetLogger2()),
			statusProvider: statusProvider,
		}
		documentService.graphRag.SetOnEvents(documentService.initOnEvents)
	})
	return documentService
}

func NewDocumentServiceDefault() *DocumentService {
	return NewDocumentService(outside.NewImportStatusProvider())
}

func (s *DocumentService) initOnEvents(e *storage.DocEvents) {
	e.OnStartInsert = s.onStartInsert
	e.OnDoneInsert = s.onDoneInsert
}

func (s *DocumentService) onStartInsert(ctx context.Context, ragDoc *entity.Document) {
	_ = s.updateState(ctx, s.current, model.DocumentState_Importing, "导入中...")
}

func (s *DocumentService) onDoneInsert(ctx context.Context, ragDoc *entity.Document, err error) {
	if err == nil {
		_ = s.updateState(ctx, s.current, model.DocumentState_Succee, "导入成功")
	} else {
		_ = s.updateState(ctx, s.current, model.DocumentState_Succee, err.Error())
	}
}

func (s *DocumentService) onChunkCount(ctx context.Context, ragDoc *entity.Document, chunkCount int) {

}

func (s *DocumentService) onStartExtractEntities(ctx context.Context, ragDoc *entity.Document, source *storage.Source) {
}

func (s *DocumentService) onDoneExtractEntities(ctx context.Context, ragDoc *entity.Document, source *storage.Source, err error) {

}

func (s *DocumentService) Scan(userCtx context.Context, cmd *command.DocumentScanCommand) {
	go func() {
		s.scan(userCtx, cmd.Data.TenantId)
	}()
}

// scan
// @Description:
// @receiver s
// @param ctx
// @param tenantId
// @param caseId
func (s *DocumentService) scan(ctx context.Context, tenantId string) {
	isLockScan := false
	env := env.GetEnv()
	gp.Try(func() error {
		caseId := ""
		for {
			_, _ = s.unlockScan(ctx, tenantId, caseId)

			isLockVal, err := s.lockScan(ctx, tenantId, caseId)
			if err != nil {
				return err
			}
			isLockScan = isLockVal
			if isLockScan == false {
				continue
			}
			list, err := s.getDocumentByPending(ctx)
			if err != nil {
				return err
			}
			if len(list) == 0 {
				continue
			}
			time.Sleep(time.Second * 10)
			for _, doc := range list {
				s.current = doc
				err = s.updateState(ctx, doc, model.DocumentState_Importing, "")
				if err != nil {
					msg := fmt.Sprintf("updateState() error:%s", doc.FsKey)
					logs.Errorfmt(ctx, msg)
				}
				fs, ok := env.Fsm.GetFs(doc.FsKey)
				if !ok {
					msg := fmt.Sprintf("fskey=%s not found", doc.FsKey)
					_ = s.updateState(ctx, doc, model.DocumentState_Failure, msg)
					logs.Errorfmt(ctx, msg)
					continue
				}

				fileName := fmt.Sprintf("%s/%s", doc.FilePath, doc.FileName)
				text, err := s.docExtract.Extract(fs, fileName)
				if err != nil {
					msg := fmt.Sprintf("fskey=%s, fileName=%s, extract error:%s", doc.FsKey, doc.FileName, err.Error())
					_ = s.updateState(ctx, doc, model.DocumentState_Succee, msg)
					logs.Errorfmt(ctx, msg)
					continue
				}
				ragDoc := &entity.Document{
					Id:       doc.Id,
					Text:     text,
					FileName: doc.FileName,
					TenantId: doc.TenantId,
					CaseId:   doc.CaseId,
				}
				_, err = s.graphRag.IngestDocument(ctx, ragDoc)
				if err != nil {
					logs.ErrorMsg(ctx, "rag scan UpdateStatus", err.Error())
					_ = s.updateState(ctx, doc, model.DocumentState_Failure, err.Error())
					continue
				}
				_ = s.updateState(ctx, doc, model.DocumentState_Succee, "")
			}
		}
	}).Catch(func(err error) {
		logs.Errorfmt(ctx, "rag scan err:%s", err.Error())
	}).Finally(func() {
		if isLockScan {
			_, _ = s.unlockScan(ctx, tenantId, "")
		}
	})
}

func (s *DocumentService) getDocumentByPending(ctx context.Context) ([]*model.Document, error) {
	sqlBuild := rsql.NewBuilder().Eq("state", model.DocumentState_Pending.String())
	findQuery := &store.FindPagingQueryRequest{
		Filter:   sqlBuild.Build(),
		PageNum:  0,
		PageSize: 1,
	}
	findResult := s.dao.FindPaging(ctx, findQuery)
	return findResult.GetData(), findResult.GetError()
}

func (s *DocumentService) updateState(ctx context.Context, doc *model.Document, state model.DocumentState, msg string) error {
	return gp.Try(func() error {
		doc.State = state
		doc.Message = msg
		callOptions := idao.NewCallOptions()
		callOptions.SetUpdateFields([]string{"state", "message"})
		_ = s.dao.Update(ctx, doc, callOptions)
		_ = s.statusProvider.UpdateStatus(ctx, doc, state.String(), msg)
		return nil
	}).Error
}

// Ingests 提取文档知识
func (s *DocumentService) Ingests(ctx context.Context, document []*model.Document, streams ...func(txt string)) []error {
	var docList []*entity.Document
	for _, item := range document {
		doc := &entity.Document{
			Id:       item.Id,
			FileName: item.FileName,
			TenantId: item.TenantId,
			CaseId:   item.CaseId,
		}
		docList = append(docList, doc)
	}
	_, _, errs := s.graphRag.IngestDocuments(ctx, docList)
	return errs
}

func (s *DocumentService) Create(ctx context.Context, cmd *command.DocumentCreateCommand) error {
	doc := newDocumentWithCreateCommand(ctx, cmd)
	err := s.create(ctx, doc)
	return err
}

// Create 创建文档
func (s *DocumentService) create(ctx context.Context, entity *model.Document, opts ...idao.CallOptions) error {
	if entity == nil {
		return errors.New("DocumentService.Create(entity): entity is nil")
	}
	tenantId, _ := appctx.GetTenantId(ctx)

	errs := errors.NewVerifyError()
	if entity.FileId == "" {
		errs.AppendField("fileId", "不能为空")
	}
	if entity.FileName == "" {
		errs.AppendField("fileName", "不能为空")
	}
	if entity.CaseId == "" {
		errs.AppendField("caseId", "不能为空")
	}
	if tenantId == "" {
		errs.AppendField("tenantId", "不能为空")
	}
	entity.State = model.DocumentState_Pending
	if errs.HasError() {
		return errs
	}
	entity.TenantId = tenantId
	if ok := s.docExtract.IsSupport(entity.FileName); !ok {
		return errors.New("不支持对%s文件进行解析, 文件类型不正确。", entity.FileName)
	}
	s.dao.Create(ctx, entity, opts...)
	s.statusProvider.UpdateStatus(ctx, entity, entity.State.String(), "")
	return nil
}

// Delete 删除文档
func (s *DocumentService) Delete(ctx context.Context, cmd *command.DocumentDeleteCommand, opts ...idao.CallOptions) error {
	return s.delete(ctx, cmd.Data.CaseId, cmd.Data.Id)
}

// Delete 删除文档
func (s *DocumentService) delete(ctx context.Context, caseId, id string, opts ...idao.CallOptions) error {
	errs := errors.NewVerifyError()
	tenantId, _ := appctx.GetTenantId(ctx)
	if tenantId == "" {
		errs.AppendField("tenantId", "不能为空")
	}
	if id == "" {
		errs.AppendField("id", "不能为空")
	}
	if caseId == "" {
		errs.AppendField("caseId", "不能为空")
	}
	if errs.HasError() {
		return errs
	}
	s.dao.DeleteByRSQL(ctx, fmt.Sprintf(`tenant_id=="%s" and case_id="%s" and id=="%s"`, tenantId, caseId, id), opts...)
	return s.graphRag.DeleteDoc(ctx, tenantId, caseId, id)
}

func (s *DocumentService) lockScan(ctx context.Context, tenantId, caseId string) (bool, error) {
	storeName := "lock-store"
	lockOwner := "rag"
	resourceId := fmt.Sprintf("scan-%s-%s", tenantId, caseId)
	dapr := env.GetEnv().Dapr
	if dapr == nil {
		panic("dapr is nil")
	}
	cli := dapr.GetClient()
	resp, err := cli.TryLockAlpha1(ctx, storeName, &client.LockRequest{
		LockOwner:       lockOwner,
		ResourceID:      resourceId,
		ExpiryInSeconds: 60 * 10,
	})
	if err != nil {
		panic(err)
	}
	return resp.Success, err
}

func (s *DocumentService) unlockScan(ctx context.Context, tenantId, caseId string) (bool, error) {
	storeName := "lock-store"
	lockOwner := "rag"
	resourceId := fmt.Sprintf("scan-%s-%s", tenantId, caseId)
	cli := dapr.GetDaprClient()
	_, err := cli.UnlockAlpha1(ctx, storeName, &client.UnlockRequest{
		LockOwner:  lockOwner,
		ResourceID: resourceId,
	})
	if err != nil {
		return false, err
	}
	return true, err
}

func newDocumentWithCreateCommand(ctx context.Context, cmd *command.DocumentCreateCommand) *model.Document {
	doc := &model.Document{
		FsKey:      cmd.Data.FsKey,
		FileId:     cmd.Data.FileId,
		FileName:   cmd.Data.FileName,
		FilePath:   cmd.Data.FilePath,
		SourceType: cmd.Data.SourceType,
		SourceApp:  cmd.Data.SourceApp,
		SourceId:   cmd.Data.SourceId,
		State:      model.DocumentState_Pending,
		Message:    "创建",
	}
	tenantId := appctx.GetTenantId2(ctx)
	doc.Id = cmd.Data.Id
	doc.CaseId = cmd.Data.CaseId
	doc.TenantId = tenantId
	return doc
}
