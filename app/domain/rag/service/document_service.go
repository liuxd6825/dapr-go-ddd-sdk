package service

import (
	"context"
	"fmt"
	"github.com/dapr/go-sdk/client"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/entity"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"sync"
)

type DocumentService struct {
	*dao.DocumentDao
	graphRag *my_rag.GraphRag
}

var documentService *DocumentService
var documentServiceOnce sync.Once

func NewDocumentService() *DocumentService {
	documentServiceOnce.Do(func() {
		documentService = &DocumentService{
			DocumentDao: dao.NewDocumentDao(DBKey),
			graphRag:    NewGraphRag(),
		}
	})
	return documentService
}

func (s *DocumentService) Scan(ctx context.Context, tenantId, caseId string) {
	isLockScan := false
	gp.Try(func() error {
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
				TenantId: tenantId,
				Filter:   fmt.Sprintf(`tenant_id=="%s" and case_id=="%s" and state==0`, tenantId, caseId),
				PageNum:  1,
				PageSize: 1,
			}
			findResult := s.FindPaging(ctx, findQuery)
			if len(findResult.Data) == 0 {
				break
			}
			// 处理所有文档
			list := findResult.Data
			for _, item := range list {
				doc := &entity.Document{
					Id:       item.Id,
					FileName: item.FileName,
					TenantId: item.TenantId,
					CaseId:   item.CaseId,
				}
				_, err = s.graphRag.IngestDocument(ctx, doc)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}).Finally(func() {
		if isLockScan {
			_, _ = s.unlockScan(ctx, tenantId, caseId)
		}
	})

}

func (s *DocumentService) lockScan(ctx context.Context, tenantId, caseId string) (bool, error) {
	storeName := "lock-store"
	lockOwner := "rag"
	resourceId := fmt.Sprintf("scan-%s-%s", tenantId, caseId)
	cli := dapr.GetDaprClient()
	resp, err := cli.TryLockAlpha1(ctx, storeName, &client.LockRequest{
		LockOwner:       lockOwner,
		ResourceID:      resourceId,
		ExpiryInSeconds: 60,
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

func (s *DocumentService) IngestDocuments(ctx context.Context, document []*model.Document, streams ...func(txt string)) []error {
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

func (s *DocumentService) DeleteDocument(ctx context.Context, tenantId, caseId, docId string) error {
	return s.graphRag.DeleteDoc(ctx, tenantId, caseId, docId)
}

func (s *DocumentService) Create(ctx context.Context, entity *model.Document, opts ...*idao.CallOptions) error {
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
	s.DocumentDao.Create(ctx, entity, opts...)

	return nil
}

func (s *DocumentService) Delete(ctx context.Context, caseId, id string, opts ...*idao.CallOptions) error {
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
	s.DocumentDao.DeleteByRSQL(ctx, fmt.Sprintf(`tenant_id=="%s" and case_id="%s" and id=="%s"`, tenantId, caseId, id), opts...)
	return s.graphRag.DeleteDoc(ctx, tenantId, caseId, id)
}
