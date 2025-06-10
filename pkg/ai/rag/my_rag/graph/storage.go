package graph

import (
	"context"
)

type Storage interface {
	GetKnowledge(ctx context.Context, tenantId string, caseId string, keys []string, maxDeep int) ([]string, error)
}
