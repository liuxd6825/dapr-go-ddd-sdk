package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/pkg/mxgraph"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/restapi/request"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/service/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/service/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"sync"
)

type GraphService struct {
	graphDao *dao.GraphDao
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
	graphService := &GraphService{
		graphDao: dao.NewGraphDao(),
	}

	return graphService
}

func (s *GraphService) Save(ctx context.Context, caseId string, drawId string, saveRequest *request.SaveFileRequest) error {
	if saveRequest == nil {
		return errors.New("diff name is empty")
	}

	saveBatch := s.GetSaveBatch(caseId, saveRequest.Diff)
	if saveBatch != nil {
		s.graphDao.BatchSave(ctx, saveBatch, drawId)
	}
	return nil
}

func (s *GraphService) GetSaveBatch(caseId string, fileDiff *mxgraph.FileDiff) *model.SaveBatch {
	saveBatch := model.NewSaveBatch()
	for _, update := range fileDiff.U {
		cells := update.Cells
		if cells != nil {
			// 删除内容
			for _, cell := range cells.R {
				cell.State = mxgraph.DiffState_Remove
				if cell.IsNode() {
					s.addRemoveNode(saveBatch, caseId, cell)
				} else if cell.IsEdge() {
					s.addRemoveRelation(saveBatch, caseId, cell)
				}
			}

			// 新建内容
			for _, cell := range cells.I {
				cell.State = mxgraph.DiffState_Insert
				if cell.IsNode() {
					s.addCreateNode(saveBatch, caseId, cell)
				} else if cell.IsEdge() {
					s.addCreateRelation(saveBatch, caseId, cell)
				}
			}

			// 更新内容
			for id, cell := range cells.U {
				cell.Id = id
				cell.State = mxgraph.DiffState_Update
				if cell.IsNode() {
					s.addUpdateNode(saveBatch, caseId, cell)
				} else if cell.IsEdge() {
					s.addUpdateRelation(saveBatch, caseId, cell)
				}
			}

		}
	}
	return saveBatch
}

func (s *GraphService) addCreateNode(saveBatch *model.SaveBatch, caseId string, cell *mxgraph.DiffCell) {
	if !cell.IsNode() {
		return
	}
	node := newNode(caseId, cell)
	saveBatch.Nodes.AddCreate(node)
}

func (s *GraphService) addCreateRelation(saveBatch *model.SaveBatch, caseId string, cell *mxgraph.DiffCell) {
	if !cell.IsEdge() {
		return
	}
	items := newRelation(caseId, cell)
	for _, rel := range items {
		if rel.StartId != "" && rel.EndId != "" && rel.RelType != "" {
			saveBatch.Relations.AddCreate(rel)
		}
	}
}

func (s *GraphService) addUpdateNode(saveBatch *model.SaveBatch, caseId string, cell *mxgraph.DiffCell) {
	if !cell.IsNode() {
		return
	}
	// 没有变化 退出
	if cell.XmlValue == nil {
		return
	}
	node := newNode(caseId, cell)
	saveBatch.Nodes.AddUpdate(node)

}

func (s *GraphService) addUpdateRelation(saveBatch *model.SaveBatch, caseId string, cell *mxgraph.DiffCell) {
	if !cell.IsEdge() {
		return
	}

	if cell.Extend.Type == "edgeLabel" {
		if cell.Value != nil && *cell.Value != cell.Extend.OldValue {
			// 是关系类型修改，删除旧关系
			remove := model.NewRelation()
			remove.Id = cell.Extend.ParentId + "-" + cell.Id
			remove.CaseId = caseId
			saveBatch.Relations.AddRemove(remove)
			// 创建新关系
			create := model.NewRelation()
			create.Id = cell.Extend.ParentId + "-" + cell.Id
			create.CaseId = caseId
			create.RelType = cell.GetRelType()
			create.StartId = cell.GetSourceId()
			create.EndId = cell.GetTargetId()
			create.Name = *cell.Value
			saveBatch.Relations.AddCreate(create)
		}

	} else if cell.Extend.Type == "edge" {
		if cell.Source == nil && cell.Target == nil {
			return
		}
		for _, label := range cell.Extend.Labels {
			// 是关系类型修改，删除旧关系
			remove := model.NewRelation()
			remove.Id = cell.Id + "-" + label.Id
			saveBatch.Relations.AddRemove(remove)

			// 创建新关系
			rel := model.NewRelation()
			rel.Id = cell.Id + "-" + label.Id
			rel.CaseId = caseId
			rel.RelType = label.Value
			rel.StartId = cell.GetSourceId()
			rel.EndId = cell.GetTargetId()
			saveBatch.Relations.AddCreate(rel)
		}
	}
}

func (s *GraphService) addRemoveNode(saveBatch *model.SaveBatch, caseId string, cell *mxgraph.DiffCell) {
	if !cell.IsNode() {
		return
	}
	node := newNode(caseId, cell)
	saveBatch.Nodes.AddRemove(node)
}

func (s *GraphService) addRemoveRelation(saveBatch *model.SaveBatch, caseId string, cell *mxgraph.DiffCell) {
	if !cell.IsEdge() {
		return
	}
	items := newRelation(caseId, cell)
	for _, rel := range items {
		saveBatch.Relations.AddRemove(rel)
	}
}

func newNode(caseId string, cell *mxgraph.DiffCell) *model.Node {
	node := model.NewNode()
	node.Id = cell.Id
	node.CaseId = caseId
	node.Name = cell.GetNodeName()
	node.Label = cell.GetNodeLabel()
	return node
}

func newRelation(caseId string, cell *mxgraph.DiffCell) []*model.Relation {
	var items []*model.Relation
	if cell.Extend.Type == "edgeLabel" {
		if cell.Value != nil && *cell.Value != cell.Extend.OldValue {
			// 是关系类型修改，删除旧关系
			rel := model.NewRelation()
			rel.Id = cell.Extend.ParentId + "-" + cell.Id
			rel.CaseId = caseId
			items = append(items, rel)
		}

		rel := model.NewRelation()
		rel.Id = cell.Extend.ParentId + "-" + cell.Id
		rel.CaseId = caseId
		rel.RelType = cell.GetRelType()
		rel.StartId = cell.GetSourceId()
		rel.EndId = cell.GetTargetId()
		items = append(items, rel)

	} else if cell.Extend.Type == "edge" {
		if cell.Value != nil && cell.Extend.OldValue != *cell.Value {
			// 是关系类型修改，删除旧关系
			remove := model.NewRelation()
			remove.Id = cell.Id
			remove.CaseId = caseId
			items = append(items, remove)

			rel := model.NewRelation()
			rel.Id = cell.Id
			rel.CaseId = caseId
			rel.RelType = *cell.Value
			rel.StartId = cell.Extend.SourceId
			rel.EndId = cell.Extend.TargetId
			items = append(items, rel)
		}

		for _, label := range cell.Extend.Labels {
			rel := model.NewRelation()
			rel.Id = cell.Id + "-" + label.Id
			rel.CaseId = caseId
			rel.RelType = label.Value
			rel.StartId = cell.Extend.SourceId
			rel.EndId = cell.Extend.TargetId
			items = append(items, rel)
		}
	}
	return items
}
