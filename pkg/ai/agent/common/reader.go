package common

import (
	"github.com/cloudwego/eino/schema"
	"io"
	"strings"
)

func Reader(reader *schema.StreamReader[*schema.Message]) (*strings.Builder, error) {
	content := &strings.Builder{}
	defer reader.Close()
	for {
		chunk, err := reader.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			// 错误处理
			return nil, err
		}
		// 响应片段处理
		content.WriteString(chunk.Content)
	}
	return content, nil
}
