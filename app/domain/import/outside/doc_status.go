package outside

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/service/interfaces"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"
)

type DocumentMetaService struct {
	service *service.DocumentMetaService
}

func NewImportStatusProvider(dbKey string) interfaces.ImportStatusProvider {
	return NewDocStatusProvider(dbKey)
}

func NewDocStatusProvider(dbKey string) *DocumentMetaService {
	return &DocumentMetaService{
		service: service.NewDocumentMetaService(dbKey),
	}
}

func (d *DocumentMetaService) UpdateStatus(ctx context.Context, doc interfaces.ImportDoc, status string, statusMsg string) error {
	cmd := &command.DocumentMetaSaveStatusBySourceType{}
	cmd.CommandId = idutils.NewId()
	cmd.Data = model.DocumentMeta{}
	cmd.Data.Id = idutils.NewId()
	cmd.Data.CaseId = doc.GetCaseId()
	cmd.Data.DocumentId = doc.GetSourceId()
	cmd.Data.SourceType = doc.GetSourceType()
	cmd.Data.Source = doc.GetId()
	cmd.Data.Name = "status"
	cmd.Data.Value = status
	return d.service.SaveStatusBySourceType(ctx, cmd)
}
