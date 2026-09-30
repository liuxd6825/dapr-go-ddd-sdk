package dao

import (
	"context"

	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/draw/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/graph"
)

type IGraphDao interface {
	BatchSave(ctx context.Context, batch *model2.SaveBatch, drawId string)
	FindGraphByDrawId(ctx context.Context, caseId, drawId string) *graph.GraphView
	FindInCaseByNames(ctx context.Context, caseId string, names []string) *graph.GraphView
	AddNodeList(target map[string]*model2.NodeView, source []any)
	AddRelationList(target map[string]*model2.RelationView, source []any)
}
