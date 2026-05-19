package context

import (
	"testing"
)

func TestNewContextItem(t *testing.T) {
	item := NewContextItem(T0_SystemPrompt, "你是一个知识助手", "system-1")
	if item.Priority != T0_SystemPrompt {
		t.Errorf("expected T0_SystemPrompt, got %v", item.Priority)
	}
	if item.Content != "你是一个知识助手" {
		t.Errorf("expected content '你是一个知识助手', got %s", item.Content)
	}
	if item.ID != "system-1" {
		t.Errorf("expected ID 'system-1', got %s", item.ID)
	}
}

func TestPriorityString(t *testing.T) {
	tests := []struct {
		priority Priority
		expected string
	}{
		{T0_SystemPrompt, "T0_SystemPrompt"},
		{T1_RetrievedKnowledge, "T1_RetrievedKnowledge"},
		{T2_RecentHistory, "T2_RecentHistory"},
		{T3_LongTermHistory, "T3_LongTermHistory"},
	}

	for _, tt := range tests {
		if tt.priority.String() != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, tt.priority.String())
		}
	}

	if T0_UserQuery.String() != "T0_UserQuery" {
		t.Errorf("expected T0_UserQuery, got %s", T0_UserQuery.String())
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.MaxTotalTokens != 8000 {
		t.Errorf("expected MaxTotalTokens 8000, got %d", cfg.MaxTotalTokens)
	}
	if cfg.RAGRatio != 0.7 {
		t.Errorf("expected RAGRatio 0.7, got %f", cfg.RAGRatio)
	}
	if cfg.MaxHistoryTurns != 5 {
		t.Errorf("expected MaxHistoryTurns 5, got %d", cfg.MaxHistoryTurns)
	}
}

func TestNewTokenBudget(t *testing.T) {
	budget, err := NewTokenBudget(8000, 500)
	if err != nil {
		t.Fatalf("failed to create TokenBudget: %v", err)
	}

	if budget.MaxTokens() != 8000 {
		t.Errorf("expected MaxTokens 8000, got %d", budget.MaxTokens())
	}
	if budget.Reserved() != 500 {
		t.Errorf("expected Reserved 500, got %d", budget.Reserved())
	}
	if budget.Available() != 7500 {
		t.Errorf("expected Available 7500, got %d", budget.Available())
	}
}

func TestTokenBudgetCountTokens(t *testing.T) {
	budget, err := NewTokenBudget(8000, 500)
	if err != nil {
		t.Fatalf("failed to create TokenBudget: %v", err)
	}

	text := "你好，世界"
	tokens := budget.CountTokens(text)
	if tokens <= 0 {
		t.Errorf("expected tokens > 0, got %d", tokens)
	}
}

func TestTokenBudgetAllocate(t *testing.T) {
	budget, err := NewTokenBudget(8000, 500)
	if err != nil {
		t.Fatalf("failed to create TokenBudget: %v", err)
	}

	ragBudget := budget.Allocate(0.7)
	expected := int(7500 * 0.7)
	if ragBudget != expected {
		t.Errorf("expected %d, got %d", expected, ragBudget)
	}
}

func TestTokenBudgetConsume(t *testing.T) {
	budget, err := NewTokenBudget(8000, 500)
	if err != nil {
		t.Fatalf("failed to create TokenBudget: %v", err)
	}

	budget.Consume(1000)
	if budget.Available() != 6500 {
		t.Errorf("expected 6500, got %d", budget.Available())
	}

	budget.Consume(10000)
	if budget.Available() != 0 {
		t.Errorf("expected 0, got %d", budget.Available())
	}
}

func TestNewPriorityManager(t *testing.T) {
	pm := NewPriorityManager()
	if pm == nil {
		t.Fatal("expected non-nil PriorityManager")
	}

	pm.AddContent(T0_SystemPrompt, "系统指令", "sys-1")
	pm.AddContent(T0_UserQuery, "用户问题", "user-1")
	pm.AddContent(T1_RetrievedKnowledge, "知识1", "rag-1")
	pm.AddContent(T2_RecentHistory, "历史1", "hist-1")
	pm.AddContent(T3_LongTermHistory, "摘要1", "sum-1")

	// Verify counts by priority
	if len(pm.GetByPriority(T0_SystemPrompt)) != 1 {
		t.Errorf("expected 1 T0_SystemPrompt item, got %d", len(pm.GetByPriority(T0_SystemPrompt)))
	}
	if len(pm.GetByPriority(T0_UserQuery)) != 1 {
		t.Errorf("expected 1 T0_UserQuery item, got %d", len(pm.GetByPriority(T0_UserQuery)))
	}

	items := pm.SortForAssembly()
	if len(items) != 5 {
		t.Errorf("expected 5 items, got %d", len(items))
	}

	if items[0].Priority != T0_SystemPrompt {
		t.Errorf("first item should be T0_SystemPrompt, got %v", items[0].Priority)
	}
	if items[len(items)-1].Priority != T0_UserQuery {
		t.Errorf("last item should be T0_UserQuery, got %v", items[len(items)-1].Priority)
	}
}

func TestPriorityManagerMerge(t *testing.T) {
	pm1 := NewPriorityManager()
	pm1.AddContent(T0_SystemPrompt, "sys-1", "id1")

	pm2 := NewPriorityManager()
	pm2.AddContent(T0_SystemPrompt, "sys-2", "id2")

	pm1.Merge(pm2)
	items := pm1.GetByPriority(T0_SystemPrompt)
	if len(items) != 2 {
		t.Errorf("expected 2 items after merge, got %d", len(items))
	}
}

func TestContentHelpers(t *testing.T) {
	userContent := NewUserContent("Hello")
	if userContent.Role != "user" {
		t.Errorf("expected role 'user', got %s", userContent.Role)
	}

	assistantContent := NewAssistantContent("Hi there")
	if assistantContent.Role != "assistant" {
		t.Errorf("expected role 'assistant', got %s", assistantContent.Role)
	}
}

func TestBuildTemplateData(t *testing.T) {
	param := &BuildParam{
		Query: "什么是GraphRAG",
		RetrievedContext: []string{"GraphRAG是一种混合检索和图结构的技术"},
		ConversationHistory: []*Content{
			{Role: "user", Content: "你好"},
			{Role: "assistant", Content: "你好，有什么可以帮助你的？"},
		},
	}

	data := BuildTemplateData(param, "用户询问GraphRAG相关信息")
	if data.CurrentQuery != "什么是GraphRAG" {
		t.Errorf("expected CurrentQuery '什么是GraphRAG', got %s", data.CurrentQuery)
	}
	if len(data.RetrievedContext) != 1 {
		t.Errorf("expected 1 retrieved context, got %d", len(data.RetrievedContext))
	}
	if len(data.History) != 2 {
		t.Errorf("expected 2 history items, got %d", len(data.History))
	}
}

func TestRenderPrompt(t *testing.T) {
	data := &TemplateData{
		SystemPrompt:     "你是一个助手",
		Summary:          "用户询问技术问题",
		History:          []*Content{{Role: "user", Content: "你好"}},
		RetrievedContext: []string{"上下文1", "上下文2"},
		CurrentQuery:     "GraphRAG是什么？",
	}

	prompt, err := RenderPrompt(data)
	if err != nil {
		t.Fatalf("failed to render prompt: %v", err)
	}

	if prompt == "" {
		t.Fatal("expected non-empty prompt")
	}

	if !contains(prompt, "你是一个助手") {
		t.Error("prompt should contain SystemPrompt")
	}
	if !contains(prompt, "GraphRAG是什么？") {
		t.Error("prompt should contain CurrentQuery")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}