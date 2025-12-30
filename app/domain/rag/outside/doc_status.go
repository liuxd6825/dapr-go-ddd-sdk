package outside

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/service/interfaces"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"
)

type DocumentMetaService struct {
	service *service.DocumentMetaService
}

func NewImportStatusProvider() interfaces.ImportStatusProvider {
	return NewDocStatusProvider()
}

func NewDocStatusProvider() *DocumentMetaService {
	return &DocumentMetaService{
		service: service.NewDocumentMetaService(),
	}
}

func (d *DocumentMetaService) UpdateStatus(ctx context.Context, doc interfaces.ImportDoc, status string, statusMsg string) error {
	if doc.GetSourceApp() == "document_service" {
		cmd := &command.DocumentMetaSaveStatusBySourceType{}
		cmd.CommandId = idutils.NewId()
		cmd.Data = model.DocumentMeta{}
		cmd.Data.Id = idutils.NewId()
		cmd.Data.CaseId = doc.GetCaseId()
		cmd.Data.DocumentId = doc.GetSourceId()
		cmd.Data.SourceType = "知识"
		cmd.Data.Source = doc.GetId()
		cmd.Data.Name = "status"
		cmd.Data.Value = status
		return d.service.SaveStatusBySourceType(ctx, cmd)
	}
	return nil
}
