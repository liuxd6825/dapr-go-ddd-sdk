package dao

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/graph"
)

type IMasterGraphDao interface {
	FindByCaseId(ctx context.Context, caseId string) *graph.GraphView
	FindById(ctx context.Context, caseId string, id string) *graph.GraphView
	FindByName(ctx context.Context, caseId string, name string) *graph.GraphView
	FindByContainName(ctx context.Context, caseId string, name string) *graph.GraphView
	FindByStartWithName(ctx context.Context, caseId string, name string) *graph.GraphView
}
