package neo4jagent

import (
	"context"
	"github.com/cloudwego/eino/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/agent/common"
	"testing"
)

func Test_newText2CypherChatModel(t *testing.T) {
	ctx := context.Background()
	chatModel, err := newText2CypherChatModel(ctx)
	if err != nil {
		t.Fatal(err)
		return
	}
	msg := []*schema.Message{}
	msg = append(msg, &schema.Message{Role: "user", Content: "将“用户名称为张三的人的关系”转换为cypher"})
	reader, err := chatModel.Stream(ctx, msg)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer reader.Close()
	content, err := common.Reader(reader)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(content.String())
}

func Test_newOllamaChatModel(t *testing.T) {
	ctx := context.Background()
	chatModel, err := newOllamaChatModel(ctx)
	if err != nil {
		t.Fatal(err)
		return
	}
	msg := []*schema.Message{}
	msg = append(msg, &schema.Message{Role: "user", Content: "将“用户名称为张三的人的关系”转换为cypher"})
	reader, err := chatModel.Stream(ctx, msg)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer reader.Close()
	content, err := common.Reader(reader)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(content.String())
}

func Test_newDeepSeekR1ForAli(t *testing.T) {
	ctx := context.Background()
	model := "deepseek-r1-distill-llama-70b"
	// model := "deepseek-r1", // 使用的模型版本

	chatModel, err := newDeepSeekR1ForAli(ctx, model)
	if err != nil {
		t.Fatal(err)
		return
	}
	msg := []*schema.Message{}
	msg = append(msg, &schema.Message{Role: "user", Content: "将“用户名称为张三的人的关系”转换为cypher"})
	reader, err := chatModel.Stream(ctx, msg)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer reader.Close()
	content, err := common.Reader(reader)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(content.String())
}
