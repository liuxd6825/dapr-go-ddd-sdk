package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/pkg/mxgraph"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/draw/dao"
	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/draw/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/draw/service/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
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

func (s *GraphService) Save(ctx context.Context, cmd command.IDrawSaveCommand) error {
	if cmd == nil {
		return errors.New("diff name is empty")
	}

	saveBatch := s.GetSaveBatch(cmd)
	if saveBatch != nil {
		s.graphDao.BatchSave(ctx, saveBatch, cmd.GetDrawId())
	}
	return nil
}

func (s *GraphService) GetSaveBatch(cmd command.IDrawSaveCommand) *model2.SaveBatch {
	saveBatch := model2.NewSaveBatch()
	for _, update := range cmd.GetDiff().U {
		cells := update.Cells
		if cells != nil {
			// 删除内容
			for _, cell := range cells.R {
				cell.State = mxgraph.DiffState_Remove
				if cell.IsNode() {
					s.addRemoveNode(saveBatch, cmd, cell)
				} else if cell.IsEdge() {
					s.addRemoveRelation(saveBatch, cmd, cell)
				}
			}

			// 新建内容
			for _, cell := range cells.I {
				cell.State = mxgraph.DiffState_Insert
				if cell.IsNode() {
					s.addCreateNode(saveBatch, cmd, cell)
				} else if cell.IsEdge() {
					s.addCreateRelation(saveBatch, cmd, cell)
				}
			}

			// 更新内容
			for id, cell := range cells.U {
				cell.Id = id
				cell.State = mxgraph.DiffState_Update
				if cell.IsNode() {
					s.addUpdateNode(saveBatch, cmd, cell)
				} else if cell.IsEdge() {
					s.addUpdateRelation(saveBatch, cmd, cell)
				}
			}

		}
	}
	return saveBatch
}

func (s *GraphService) addCreateNode(saveBatch *model2.SaveBatch, cmd command.IDrawSaveCommand, cell *mxgraph.DiffCell) {
	if !cell.IsNode() {
		return
	}
	node := newNode(cmd, cell)
	saveBatch.Nodes.AddCreate(node)
}

func (s *GraphService) addCreateRelation(saveBatch *model2.SaveBatch, cmd command.IDrawSaveCommand, cell *mxgraph.DiffCell) {
	if !cell.IsEdge() {
		return
	}
	items := newRelation(cmd, cell)
	for _, rel := range items {
		if rel.Source != "" && rel.Target != "" && rel.RelType != "" {
			saveBatch.Relations.AddCreate(rel)
		}
	}
}

func (s *GraphService) addUpdateNode(saveBatch *model2.SaveBatch, cmd command.IDrawSaveCommand, cell *mxgraph.DiffCell) {
	if !cell.IsNode() {
		return
	}
	// 没有变化 退出
	if cell.XmlValue == nil {
		return
	}
	node := newNode(cmd, cell)
	saveBatch.Nodes.AddUpdate(node)

}

func (s *GraphService) addUpdateRelation(saveBatch *model2.SaveBatch, cmd command.IDrawSaveCommand, cell *mxgraph.DiffCell) {
	if !cell.IsEdge() {
		return
	}

	if cell.Extend.Type == "edgeLabel" {
		if cell.Value != nil && *cell.Value != cell.Extend.OldValue {
			// 是关系类型修改，删除旧关系
			remove := model2.NewRelation()
			remove.Id = cell.Extend.ParentId + "-" + cell.Id
			remove.CaseId = cmd.GetCaseId()
			saveBatch.Relations.AddRemove(remove)
			// 创建新关系
			create := model2.NewRelation()
			create.Id = cell.Extend.ParentId + "-" + cell.Id
			create.CaseId = cmd.GetCaseId()
			create.RelType = cell.GetRelType()
			create.Source = cell.GetSourceId()
			create.Target = cell.GetTargetId()
			create.Name = *cell.Value
			create.SourceUrl = getSourceUrl(cmd)
			create.SourceName = getSourceName(cmd)
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
			rel.CaseId = cmd.GetCaseId()
			rel.RelType = label.Value
			rel.Source = cell.GetSourceId()
			rel.Target = cell.GetTargetId()
			rel.SourceUrl = getSourceUrl(cmd)
			rel.SourceName = getSourceName(cmd)
			saveBatch.Relations.AddCreate(rel)
		}
	}
}

func (s *GraphService) addRemoveNode(saveBatch *model2.SaveBatch, cmd command.IDrawSaveCommand, cell *mxgraph.DiffCell) {
	if !cell.IsNode() {
		return
	}
	node := newNode(cmd, cell)
	saveBatch.Nodes.AddRemove(node)
}

func (s *GraphService) addRemoveRelation(saveBatch *model2.SaveBatch, cmd command.IDrawSaveCommand, cell *mxgraph.DiffCell) {
	if !cell.IsEdge() {
		return
	}
	items := newRelation(cmd, cell)
	for _, rel := range items {
		saveBatch.Relations.AddRemove(rel)
	}
}

func newNode(cmd command.IDrawSaveCommand, cell *mxgraph.DiffCell) *model2.Node {
	node := model2.NewNode()
	node.Id = cell.Id
	node.CaseId = cmd.GetCaseId()
	node.Name = cell.GetNodeName()
	node.Type = cell.GetNodeLabel()
	node.SourceType = "draw"
	node.SourceIds = cell.Id
	node.Type = "draw"
	node.Description = node.Name
	node.SourceUrl = getSourceUrl(cmd)
	node.SourceName = getSourceName(cmd)
	return node
}

func newRelation(cmd command.IDrawSaveCommand, cell *mxgraph.DiffCell) []*model2.Relation {
	var items []*model2.Relation
	if cell.Extend.Type == "edgeLabel" {
		if cell.Value != nil && *cell.Value != cell.Extend.OldValue {
			// 是关系类型修改，删除旧关系
			rel := model2.NewRelation()
			rel.Id = cell.Extend.ParentId + "-" + cell.Id
			rel.CaseId = cmd.GetCaseId()
			items = append(items, rel)
		}

		rel := model2.NewRelation()
		rel.Id = cell.Extend.ParentId + "-" + cell.Id
		rel.CaseId = cmd.GetCaseId()
		rel.RelType = cell.GetRelType()
		rel.Source = cell.GetSourceId()
		rel.Target = cell.GetTargetId()
		rel.Keywords = []string{rel.RelType}
		rel.SourceUrl = getSourceUrl(cmd)
		rel.SourceName = getSourceName(cmd)
		items = append(items, rel)

	} else if cell.Extend.Type == "edge" {
		if cell.Value != nil && cell.Extend.OldValue != *cell.Value {
			// 是关系类型修改，删除旧关系
			remove := model2.NewRelation()
			remove.Id = cell.Id
			remove.CaseId = cmd.GetCaseId()
			items = append(items, remove)

			rel := model2.NewRelation()
			rel.Id = cell.Id
			rel.CaseId = cmd.GetCaseId()
			rel.RelType = *cell.Value
			rel.Source = cell.Extend.SourceId
			rel.Target = cell.Extend.TargetId
			rel.SourceUrl = getSourceUrl(cmd)
			rel.SourceName = getSourceName(cmd)
			items = append(items, rel)
		}

		for _, label := range cell.Extend.Labels {
			rel := model2.NewRelation()
			rel.Id = cell.Id + "-" + label.Id
			rel.CaseId = cmd.GetCaseId()
			rel.RelType = label.Value
			rel.Source = cell.Extend.SourceId
			rel.Target = cell.Extend.TargetId
			rel.SourceUrl = getSourceUrl(cmd)
			rel.SourceName = getSourceName(cmd)
			items = append(items, rel)
		}
	}
	return items
}

func getSourceUrl(cmd command.IDrawSaveCommand) string {
	if env.GetEnv().App.ProdMode {
		return fmt.Sprintf("/draw/draw.html?id=%s&case-id=%s", cmd.GetDrawId(), cmd.GetCaseId())
	}
	return fmt.Sprintf("/draw/draw.html?dev=1&id=%s&case-id=%s", cmd.GetDrawId(), cmd.GetCaseId())
}

func getSourceName(cmd command.IDrawSaveCommand) string {
	return cmd.GetFileName()
}
