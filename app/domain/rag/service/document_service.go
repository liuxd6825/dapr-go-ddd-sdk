package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/dapr/go-sdk/client"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/doc_extract"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/entity"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/storage"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"
)

type DocumentService struct {
	dao        *dao.DocumentDao
	graphRag   *my_rag.GraphRag
	docExtract *doc_extract.Extract
}

var documentService *DocumentService
var documentServiceOnce sync.Once

func NewDocumentService() *DocumentService {
	documentServiceOnce.Do(func() {
		documentService = &DocumentService{
			dao:        dao.NewDocumentDao(DBKey),
			graphRag:   NewGraphRag(),
			docExtract: doc_extract.NewExtract(logs.GetLogger2()),
		}
		documentService.graphRag.SetOnEvents(documentService.initOnEvents)
	})
	return documentService
}

func (s *DocumentService) initOnEvents(e *storage.DocEvents) {
	e.OnStartInsert = s.onStartInsert
	e.OnDoneInsert = s.onDoneInsert
}

func (s *DocumentService) onStartInsert(ctx context.Context, ragDoc *entity.Document) {
	_ = s.updateState(ctx, ragDoc.TenantId, ragDoc.CaseId, ragDoc.Id, 1, "导入中...")
}

func (s *DocumentService) onDoneInsert(ctx context.Context, ragDoc *entity.Document, err error) {
	if err == nil {
		_ = s.updateState(ctx, ragDoc.TenantId, ragDoc.CaseId, ragDoc.Id, 100, "导入成功")
	} else {
		_ = s.updateState(ctx, ragDoc.TenantId, ragDoc.CaseId, ragDoc.Id, -1, err.Error())
	}
}

func (s *DocumentService) onChunkCount(ctx context.Context, ragDoc *entity.Document, chunkCount int) {

}

func (s *DocumentService) onStartExtractEntities(ctx context.Context, ragDoc *entity.Document, source *storage.Source) {
}

func (s *DocumentService) onDoneExtractEntities(ctx context.Context, ragDoc *entity.Document, source *storage.Source, err error) {

}

func (s *DocumentService) Scan(ctx context.Context, cmd *command.DocumentScanCommand) {
	go func() {
		s.scan(ctx, cmd.Data.TenantId, cmd.Data.CaseId)
	}()
}

func (s *DocumentService) scan(ctx context.Context, tenantId, caseId string) {
	isLockScan := false
	gp.Try(func() error {
		_, _ = s.unlockScan(ctx, tenantId, caseId)
		isLockVal, err := s.lockScan(ctx, tenantId, caseId)
		if err != nil {
			return err
		}
		isLockScan = isLockVal
		if isLockScan == false {
			return nil
		}

		for {
			findQuery := &store.FindPagingQueryRequest{
				Filter:   fmt.Sprintf("case_id=='%s' and state==0", caseId),
				PageNum:  0,
				PageSize: 1,
			}
			findResult := s.dao.FindPaging(ctx, findQuery)
			if findResult.GetError() != nil {
				return findResult.GetError()
			}
			if len(findResult.GetData()) == 0 {
				break
			}
			// 处理所有文档
			list := findResult.GetData()
			env := env.GetEnv()
			for _, doc := range list {
				fs, ok := env.Fsm.GetFs(doc.FsKey)
				if !ok {
					msg := fmt.Sprintf("fskey=%s not found", doc.FsKey)
					_ = s.updateState(ctx, doc.TenantId, doc.CaseId, doc.Id, -1, msg)
					logs.Errorfmt(ctx, msg)
					continue
				}
				text, err := s.docExtract.Extract(fs, doc.Path+doc.FileName)
				if err != nil {
					msg := fmt.Sprintf("fskey=%s, fileName=%s, extract error:%s", doc.FsKey, doc.FileName, err.Error())
					_ = s.updateState(ctx, doc.TenantId, doc.CaseId, doc.Id, -1, msg)
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
					return err
				}
			}
		}
		return nil
	}).Catch(func(err error) {
		logs.Errorfmt(ctx, "rag scan err:%s", err.Error())
	}).Finally(func() {
		if isLockScan {
			_, _ = s.unlockScan(ctx, tenantId, caseId)
		}
	})
}

func (s *DocumentService) updateState(ctx context.Context, tenantId string, caseId, id string, state int, msg string) error {
	doc := &model.Document{}
	doc.Id = id
	doc.TenantId = tenantId
	doc.CaseId = caseId
	doc.State = state
	doc.Message = msg
	callOptions := idao.NewCallOptions()
	callOptions.SetUpdateFields([]string{"state", "message"})
	_ = s.dao.Update(ctx, doc, callOptions)
	return nil
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
	if err == nil {
		s.scan(ctx, doc.TenantId, doc.CaseId)
	}
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
	if entity.State != 0 {
		errs.AppendField("state", "状态不正确")
	}
	if errs.HasError() {
		return errs
	}
	entity.TenantId = tenantId
	if ok := s.docExtract.IsSupport(entity.FileName); !ok {
		return errors.New("不支持对%s文件进行解析, 文件类型不正确。", entity.FileName)
	}
	s.dao.Create(ctx, entity, opts...)
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
		FsKey:    cmd.Data.FsKey,
		FileId:   cmd.Data.FileId,
		FileName: cmd.Data.FileName,
		Path:     cmd.Data.Path,
		State:    0,
		Message:  "创建",
	}
	tenantId := appctx.GetTenantId2(ctx)
	doc.Id = idutils.NewId()
	doc.CaseId = cmd.Data.CaseId
	doc.TenantId = tenantId
	return doc
}

func newDocument(cmdData *command.DocumentData) *model.Document {
	return &model.Document{
		FileId:     cmdData.FileId,
		FileName:   cmdData.FileName,
		State:      cmdData.State,
		ChunkCount: cmdData.ChunkCount,
		DoneChunk:  cmdData.DoneChunk,
		StartTime:  cmdData.StartTime,
		EndTime:    cmdData.EndTime,
		Message:    cmdData.Message,
	}
}
