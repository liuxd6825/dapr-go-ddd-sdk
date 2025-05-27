package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/drawio/pkg/mxgraph"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/drawio/restapi/request"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/drawio/service/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"sync"
)

type GraphService struct {
}

var _graphService *GraphService
var _graphServiceOnce sync.Once

func NewGraphService() *GraphService {
	_graphServiceOnce.Do(func() {
		_graphService = newGraphService()
	})
	return _graphService
}

func newGraphService() *GraphService {
	graphService := &GraphService{}
	return graphService
}

func (s *GraphService) Save(ctx context.Context, saveRequest *request.SaveFileRequest) error {
	if saveRequest == nil {
		return errors.New("diff name is empty")
	}

	drawioFile, err := mxgraph.NewDrawioFile(saveRequest.XML)
	if err != nil {
		return err
	}

	s.GetSaveNodes(drawioFile, saveRequest.Diff)
	return nil
}

func (s *GraphService) GetSaveNodes(drawioFile *mxgraph.DrawioFile, fileDiff *mxgraph.FileDiff) *model.SaveBatch {
	saveBatch := model.NewSaveBatch()
	for _, update := range fileDiff.U {
		cells := update.Cells
		if cells != nil {
			// 删除内容
			rItems := cells.R
			for _, cellId := range rItems {
				println(cellId)
			}

			// 新建内容
			nItems := cells.I
			for _, cell := range nItems {
				if cell.IsNode() {
				} else if cell.IsEdge() {
				}
			}

			// 更新内容
			uItems := cells.U
			for id, cell := range uItems {
				cell.Id = id
				if cell.IsNode() {
				} else if cell.IsEdge() {
				}
			}
		}
	}
	return saveBatch
}
