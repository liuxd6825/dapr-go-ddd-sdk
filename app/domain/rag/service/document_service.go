package service

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"sync"
)

type DocumentService struct {
	*dao.DocumentDao
	rag *RagService
}

var documentService *DocumentService
var documentServiceOnce sync.Once

func NewDocumentService() *DocumentService {
	documentServiceOnce.Do(func() {
		documentService = &DocumentService{
			DocumentDao: dao.NewDocumentDao(DBKey),
			rag:         NewRagService(),
		}
	})
	return documentService
}

func (s *DocumentService) Scan(ctx context.Context, tenantId, caseId string) []error {
	list := s.FindByRSQL(ctx, fmt.Sprintf(`tenant_id=="%s" and case_id=="%s" and state==0`, tenantId, caseId))
	err := ragService.IngestDocuments(ctx, list)
	return err
}

func (s *DocumentService) Create(ctx context.Context, entity *model.Document, opts ...*idao.CallOptions) error {
	if entity == nil {
		return errors.New("DocumentService.Create(entity): entity is nil")
	}
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
	if entity.TenantId == "" {
		errs.AppendField("tenantId", "不能为空")
	}
	if entity.State == "" {
		errs.AppendField("state", "不能为空")
	}
	if errs.HasError() {
		return errs
	}
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

	s.DocumentDao.DeleteByRSQL(ctx, fmt.Sprintf(`case_id="%s" and id=="%s"`, caseId, id), opts...)
	return s.rag.DeleteDocument(ctx, tenantId, caseId, id)
}
