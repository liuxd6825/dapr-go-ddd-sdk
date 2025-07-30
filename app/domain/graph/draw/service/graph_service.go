package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/pkg/mxgraph"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/restapi/request"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/draw/dao"
	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/draw/model"
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
		graphDao: dao.NewGraphDao(service.Neo4jDBKey),
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

func (s *GraphService) GetSaveBatch(caseId string, fileDiff *mxgraph.FileDiff) *model2.SaveBatch {
	saveBatch := model2.NewSaveBatch()
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

func (s *GraphService) addCreateNode(saveBatch *model2.SaveBatch, caseId string, cell *mxgraph.DiffCell) {
	if !cell.IsNode() {
		return
	}
	node := newNode(caseId, cell)
	saveBatch.Nodes.AddCreate(node)
}

func (s *GraphService) addCreateRelation(saveBatch *model2.SaveBatch, caseId string, cell *mxgraph.DiffCell) {
	if !cell.IsEdge() {
		return
	}
	items := newRelation(caseId, cell)
	for _, rel := range items {
		if rel.Source != "" && rel.Target != "" && rel.RelType != "" {
			saveBatch.Relations.AddCreate(rel)
		}
	}
}

func (s *GraphService) addUpdateNode(saveBatch *model2.SaveBatch, caseId string, cell *mxgraph.DiffCell) {
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

func (s *GraphService) addUpdateRelation(saveBatch *model2.SaveBatch, caseId string, cell *mxgraph.DiffCell) {
	if !cell.IsEdge() {
		return
	}

	if cell.Extend.Type == "edgeLabel" {
		if cell.Value != nil && *cell.Value != cell.Extend.OldValue {
			// 是关系类型修改，删除旧关系
			remove := model2.NewRelation()
			remove.Id = cell.Extend.ParentId + "-" + cell.Id
			remove.CaseId = caseId
			saveBatch.Relations.AddRemove(remove)
			// 创建新关系
			create := model2.NewRelation()
			create.Id = cell.Extend.ParentId + "-" + cell.Id
			create.CaseId = caseId
			create.RelType = cell.GetRelType()
			create.Source = cell.GetSourceId()
			create.Target = cell.GetTargetId()
			create.Name = *cell.Value
			saveBatch.Relations.AddCreate(create)
		}

	} else if cell.Extend.Type == "edge" {
		if cell.Source == nil && cell.Target == nil {
			return
		}
		for _, label := range cell.Extend.Labels {
			// 是关系类型修改，删除旧关系
			remove := model2.NewRelation()
			remove.Id = cell.Id + "-" + label.Id
			saveBatch.Relations.AddRemove(remove)

			// 创建新关系
			rel := model2.NewRelation()
			rel.Id = cell.Id + "-" + label.Id
			rel.CaseId = caseId
			rel.RelType = label.Value
			rel.Source = cell.GetSourceId()
			rel.Target = cell.GetTargetId()
			saveBatch.Relations.AddCreate(rel)
		}
	}
}

func (s *GraphService) addRemoveNode(saveBatch *model2.SaveBatch, caseId string, cell *mxgraph.DiffCell) {
	if !cell.IsNode() {
		return
	}
	node := newNode(caseId, cell)
	saveBatch.Nodes.AddRemove(node)
}

func (s *GraphService) addRemoveRelation(saveBatch *model2.SaveBatch, caseId string, cell *mxgraph.DiffCell) {
	if !cell.IsEdge() {
		return
	}
	items := newRelation(caseId, cell)
	for _, rel := range items {
		saveBatch.Relations.AddRemove(rel)
	}
}

func newNode(caseId string, cell *mxgraph.DiffCell) *model2.Node {
	node := model2.NewNode()
	node.Id = cell.Id
	node.CaseId = caseId
	node.Name = cell.GetNodeName()
	node.Type = cell.GetNodeLabel()
	node.SourceType = "draw"
	node.SourceIds = cell.Id
	node.Description = ""
	return node
}

func newRelation(caseId string, cell *mxgraph.DiffCell) []*model2.Relation {
	var items []*model2.Relation
	if cell.Extend.Type == "edgeLabel" {
		if cell.Value != nil && *cell.Value != cell.Extend.OldValue {
			// 是关系类型修改，删除旧关系
			rel := model2.NewRelation()
			rel.Id = cell.Extend.ParentId + "-" + cell.Id
			rel.CaseId = caseId
			items = append(items, rel)
		}

		rel := model2.NewRelation()
		rel.Id = cell.Extend.ParentId + "-" + cell.Id
		rel.CaseId = caseId
		rel.RelType = cell.GetRelType()
		rel.Source = cell.GetSourceId()
		rel.Target = cell.GetTargetId()
		rel.Keywords = []string{rel.RelType}
		items = append(items, rel)

	} else if cell.Extend.Type == "edge" {
		if cell.Value != nil && cell.Extend.OldValue != *cell.Value {
			// 是关系类型修改，删除旧关系
			remove := model2.NewRelation()
			remove.Id = cell.Id
			remove.CaseId = caseId
			items = append(items, remove)

			rel := model2.NewRelation()
			rel.Id = cell.Id
			rel.CaseId = caseId
			rel.RelType = *cell.Value
			rel.Source = cell.Extend.SourceId
			rel.Target = cell.Extend.TargetId
			items = append(items, rel)
		}

		for _, label := range cell.Extend.Labels {
			rel := model2.NewRelation()
			rel.Id = cell.Id + "-" + label.Id
			rel.CaseId = caseId
			rel.RelType = label.Value
			rel.Source = cell.Extend.SourceId
			rel.Target = cell.Extend.TargetId
			items = append(items, rel)
		}
	}
	return items
}
