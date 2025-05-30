package request

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/pkg/mxgraph"

type SaveFileRequest struct {
	XML  string            `json:"xml"`
	Diff *mxgraph.FileDiff `json:"diff"`
}
