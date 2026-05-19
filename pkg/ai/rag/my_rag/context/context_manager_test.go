package context

import (
	"context"
	"testing"
)

func TestTruncator_NewTruncator(t *testing.T) {
	budget, _ := NewTokenBudget(8000, 500)
	truncator := NewTruncator(budget, 5)

	if truncator.GetMaxTurns() != 5 {
		t.Errorf("expected maxTurns 5, got %d", truncator.GetMaxTurns())
	}

	truncator.SetMaxTurns(10)
	if truncator.GetMaxTurns() != 10 {
		t.Errorf("expected maxTurns 10, got %d", truncator.GetMaxTurns())
	}
}

func TestTruncator_TruncateEmptyHistory(t *testing.T) {
	budget, _ := NewTokenBudget(8000, 500)
	truncator := NewTruncator(budget, 5)

	result := truncator.TruncateConversationHistory([]*Content{})
	if len(result) != 0 {
		t.Errorf("expected 0 items, got %d", len(result))
	}
}

func TestTruncator_TruncateBasic(t *testing.T) {
	budget, _ := NewTokenBudget(8000, 500)
	truncator := NewTruncator(budget, 3)

	history := []*Content{
		{Role: "user", Content: "你好"},
		{Role: "assistant", Content: "你好，有什么可以帮助你的？"},
		{Role: "user", Content: "什么是GraphRAG？"},
		{Role: "assistant", Content: "GraphRAG是一种混合检索和图结构的技术。"},
		{Role: "user", Content: "它和传统RAG有什么区别？"},
		{Role: "assistant", Content: "传统RAG只依赖向量检索，而GraphRAG还利用了知识图谱。"},
	}

	result := truncator.TruncateConversationHistory(history)

	if len(result) == 0 {
		t.Fatal("expected non-empty result")
	}

	if result[0].Role != "user" {
		t.Errorf("first item should be user role, got %s", result[0].Role)
	}
}

func TestTruncator_TruncateByMaxTurns(t *testing.T) {
	budget, _ := NewTokenBudget(8000, 500)
	truncator := NewTruncator(budget, 2)

	history := []*Content{
		{Role: "user", Content: "第1轮用户"},
		{Role: "assistant", Content: "第1轮助手"},
		{Role: "user", Content: "第2轮用户"},
		{Role: "assistant", Content: "第2轮助手"},
		{Role: "user", Content: "第3轮用户"},
		{Role: "assistant", Content: "第3轮助手"},
	}

	result := truncator.TruncateConversationHistory(history)

	if len(result) > 4 {
		t.Errorf("expected at most 4 items (2 turns), got %d", len(result))
	}
}

func TestTruncator_TruncateByTokens(t *testing.T) {
	budget, _ := NewTokenBudget(8000, 500)
	truncator := NewTruncator(budget, 5)

	items := []*ContextItem{
		{ID: "1", Content: "这是第一段内容", TokenSize: 10},
		{ID: "2", Content: "这是第二段内容", TokenSize: 20},
		{ID: "3", Content: "这是第三段内容", TokenSize: 15},
		{ID: "4", Content: "这是第四段内容", TokenSize: 30},
	}

	result := truncator.TruncateByTokens(items, 35)

	total := 0
	for _, item := range result {
		total += item.TokenSize
	}

	if total > 35 {
		t.Errorf("expected total tokens <= 35, got %d", total)
	}
}

func TestContextManager_NewContextManager(t *testing.T) {
	cfg := DefaultConfig()
	cm, err := NewContextManager(cfg)
	if err != nil {
		t.Fatalf("failed to create ContextManager: %v", err)
	}

	if cm.GetConfig() != cfg {
		t.Error("config should match")
	}

	if cm.GetBudget() == nil {
		t.Error("budget should not be nil")
	}

	if cm.GetTruncator() == nil {
		t.Error("truncator should not be nil")
	}
}

func TestContextManager_BuildBasic(t *testing.T) {
	cfg := DefaultConfig()
	cm, err := NewContextManager(cfg)
	if err != nil {
		t.Fatalf("failed to create ContextManager: %v", err)
	}

	param := &BuildParam{
		Query:            "GraphRAG是什么？",
		RetrievedContext: []string{"GraphRAG是一种混合检索技术"},
	}

	result, err := cm.Build(context.Background(), param)
	if err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	if result.FinalPrompt == "" {
		t.Error("expected non-empty prompt")
	}

	if result.UsedTokens == 0 {
		t.Error("expected used tokens > 0")
	}

	if len(result.ContextItems) == 0 {
		t.Error("expected non-empty context items")
	}

	if result.ContextItems[0].Priority != T0_SystemPrompt {
		t.Error("first item should be system prompt")
	}

	lastItem := result.ContextItems[len(result.ContextItems)-1]
	if lastItem.Priority != T0_UserQuery {
		t.Error("last item should be user query")
	}
}

func TestContextManager_BuildWithHistory(t *testing.T) {
	cfg := DefaultConfig()
	cm, err := NewContextManager(cfg)
	if err != nil {
		t.Fatalf("failed to create ContextManager: %v", err)
	}

	param := &BuildParam{
		Query:   "它和传统RAG有什么区别？",
		RetrievedContext: []string{"GraphRAG利用知识图谱增强检索效果"},
		ConversationHistory: []*Content{
			{Role: "user", Content: "你好"},
			{Role: "assistant", Content: "你好，有什么可以帮助你的？"},
			{Role: "user", Content: "什么是GraphRAG？"},
			{Role: "assistant", Content: "GraphRAG是一种混合检索技术"},
		},
	}

	result, err := cm.Build(context.Background(), param)
	if err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	if result.FinalPrompt == "" {
		t.Error("expected non-empty prompt")
	}

	hasHistory := false
	for _, item := range result.ContextItems {
		if item.Priority == T2_RecentHistory {
			hasHistory = true
			break
		}
	}
	if !hasHistory {
		t.Error("expected at least one history item")
	}
}

func TestContextManager_BuildWithEmptyQuery(t *testing.T) {
	cfg := DefaultConfig()
	cm, err := NewContextManager(cfg)
	if err != nil {
		t.Fatalf("failed to create ContextManager: %v", err)
	}

	param := &BuildParam{
		Query:            "",
		RetrievedContext: []string{"GraphRAG是一种混合检索技术"},
	}

	result, err := cm.Build(context.Background(), param)
	if err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	if result.FinalPrompt == "" {
		t.Error("expected non-empty prompt")
	}
}

func TestContextManager_BuildWithEmptyContext(t *testing.T) {
	cfg := DefaultConfig()
	cm, err := NewContextManager(cfg)
	if err != nil {
		t.Fatalf("failed to create ContextManager: %v", err)
	}

	param := &BuildParam{
		Query:   "只是一个问题",
		ConversationHistory: []*Content{
			{Role: "user", Content: "你好"},
			{Role: "assistant", Content: "你好"},
		},
	}

	result, err := cm.Build(context.Background(), param)
	if err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	if result.FinalPrompt == "" {
		t.Error("expected non-empty prompt")
	}

	hasHistory := false
	for _, item := range result.ContextItems {
		if item.Priority == T2_RecentHistory {
			hasHistory = true
			break
		}
	}
	if !hasHistory {
		t.Error("expected history items")
	}
}

func TestContextManager_MultipleBuilds(t *testing.T) {
	cfg := DefaultConfig()
	cm, err := NewContextManager(cfg)
	if err != nil {
		t.Fatalf("failed to create ContextManager: %v", err)
	}

	for i := 0; i < 3; i++ {
		param := &BuildParam{
			Query:            "问题" + itoa(i),
			RetrievedContext: []string{"上下文" + itoa(i)},
		}

		result, err := cm.Build(context.Background(), param)
		if err != nil {
			t.Fatalf("build %d failed: %v", i, err)
		}

		if result.FinalPrompt == "" {
			t.Error("expected non-empty prompt")
		}
	}
}

func TestContextManager_BuildOrder(t *testing.T) {
	cfg := DefaultConfig()
	cm, err := NewContextManager(cfg)
	if err != nil {
		t.Fatalf("failed to create ContextManager: %v", err)
	}

	param := &BuildParam{
		Query:            "最终问题",
		RetrievedContext: []string{"上下文"},
		ConversationHistory: []*Content{
			{Role: "user", Content: "历史用户"},
			{Role: "assistant", Content: "历史助手"},
		},
	}

	result, err := cm.Build(context.Background(), param)
	if err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	items := result.ContextItems

	if items[0].Priority != T0_SystemPrompt {
		t.Errorf("items[0] should be T0_SystemPrompt, got %v", items[0].Priority)
	}

	lastItem := items[len(items)-1]
	if lastItem.Priority != T0_UserQuery {
		t.Error("last item should be T0_UserQuery")
	}

	hasT1 := false
	for _, item := range items {
		if item.Priority == T1_RetrievedKnowledge {
			hasT1 = true
		}
	}
	if !hasT1 {
		t.Error("expected T1_RetrievedKnowledge in items")
	}
}