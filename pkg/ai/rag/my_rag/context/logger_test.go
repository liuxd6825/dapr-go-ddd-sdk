package context

import (
	"testing"
)

func TestLogLevelString(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{LogLevelDebug, "DEBUG"},
		{LogLevelInfo, "INFO"},
		{LogLevelWarn, "WARN"},
		{LogLevelError, "ERROR"},
		{LogLevel(100), "UNKNOWN"},
	}

	for _, tt := range tests {
		if tt.level.String() != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, tt.level.String())
		}
	}
}

func TestNewContextLogger(t *testing.T) {
	logger := NewContextLogger()
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}
	if logger.Size() != 0 {
		t.Errorf("expected 0 logs, got %d", logger.Size())
	}
}

func TestContextLogger_Log(t *testing.T) {
	logger := NewContextLogger()
	logger.Clear()

	result := &BuildResult{
		FinalPrompt:    "Test prompt",
		UsedTokens:     100,
		QueryRewritten: "Test query",
		ContextItems: []*ContextItem{
			{
				ID:        "test-1",
				Priority:  T0_SystemPrompt,
				Content:   "System prompt content",
				TokenSize: 50,
			},
		},
	}

	logger.Log(result)

	if logger.Size() != 1 {
		t.Errorf("expected 1 log, got %d", logger.Size())
	}

	lastLog := logger.GetLastLog()
	if lastLog == nil {
		t.Fatal("expected non-nil last log")
	}
	if lastLog.UsedTokens != 100 {
		t.Errorf("expected UsedTokens 100, got %d", lastLog.UsedTokens)
	}
}

func TestContextLogger_LogWithLevel(t *testing.T) {
	logger := NewContextLogger()
	logger.Clear()
	logger.SetLevel(LogLevelWarn)

	result := &BuildResult{
		FinalPrompt: "Test",
		UsedTokens:  100,
		ContextItems: []*ContextItem{
			{ID: "1", Priority: T0_SystemPrompt, Content: "test", TokenSize: 10},
		},
	}

	logger.LogWithLevel(LogLevelInfo, result)

	if logger.Size() != 0 {
		t.Error("expected 0 logs since Info < Warn level")
	}

	logger.LogWithLevel(LogLevelError, result)
	if logger.Size() != 1 {
		t.Error("expected 1 log since Error >= Warn level")
	}
}

func TestContextLogger_GetLogs(t *testing.T) {
	logger := NewContextLogger()
	logger.Clear()

	for i := 0; i < 3; i++ {
		result := &BuildResult{
			UsedTokens: i * 100,
			ContextItems: []*ContextItem{
				{ID: "test", Priority: T0_SystemPrompt, Content: "test", TokenSize: 10},
			},
		}
		logger.Log(result)
	}

	logs := logger.GetLogs()
	if len(logs) != 3 {
		t.Errorf("expected 3 logs, got %d", len(logs))
	}
}

func TestContextLogger_GetLastLog(t *testing.T) {
	logger := NewContextLogger()
	logger.Clear()

	if logger.GetLastLog() != nil {
		t.Error("expected nil when no logs")
	}

	result := &BuildResult{
		UsedTokens: 100,
		ContextItems: []*ContextItem{
			{ID: "test", Priority: T0_SystemPrompt, Content: "test", TokenSize: 10},
		},
	}
	logger.Log(result)

	lastLog := logger.GetLastLog()
	if lastLog == nil {
		t.Fatal("expected non-nil last log")
	}
	if lastLog.UsedTokens != 100 {
		t.Errorf("expected UsedTokens 100, got %d", lastLog.UsedTokens)
	}
}

func TestContextLogger_Clear(t *testing.T) {
	logger := NewContextLogger()

	result := &BuildResult{
		UsedTokens: 100,
		ContextItems: []*ContextItem{
			{ID: "test", Priority: T0_SystemPrompt, Content: "test", TokenSize: 10},
		},
	}
	logger.Log(result)

	logger.Clear()
	if logger.Size() != 0 {
		t.Errorf("expected 0 logs after clear, got %d", logger.Size())
	}
}

func TestContextLogger_SetLevel(t *testing.T) {
	logger := NewContextLogger()
	logger.SetLevel(LogLevelDebug)
	logger.SetLevel(LogLevelError)
}

func TestContextLogger_SetMaxLogs(t *testing.T) {
	logger := NewContextLogger()
	logger.SetMaxLogs(5)

	result := &BuildResult{
		UsedTokens: 100,
		ContextItems: []*ContextItem{
			{ID: "test", Priority: T0_SystemPrompt, Content: "test", TokenSize: 10},
		},
	}

	for i := 0; i < 10; i++ {
		logger.Log(result)
	}

	if logger.Size() != 5 {
		t.Errorf("expected 5 logs (max), got %d", logger.Size())
	}
}

func TestContextLogger_SetRequestID(t *testing.T) {
	logger := NewContextLogger()
	logger.SetRequestID("req-123")

	result := &BuildResult{
		UsedTokens: 100,
		ContextItems: []*ContextItem{
			{ID: "test", Priority: T0_SystemPrompt, Content: "test", TokenSize: 10},
		},
	}
	logger.Log(result)

	lastLog := logger.GetLastLog()
	if lastLog.RequestID != "req-123" {
		t.Errorf("expected RequestID 'req-123', got '%s'", lastLog.RequestID)
	}
}

func TestContextLogger_EnableDisable(t *testing.T) {
	logger := NewContextLogger()
	logger.Disable()

	result := &BuildResult{
		UsedTokens: 100,
		ContextItems: []*ContextItem{
			{ID: "test", Priority: T0_SystemPrompt, Content: "test", TokenSize: 10},
		},
	}
	logger.Log(result)

	if logger.Size() != 0 {
		t.Error("expected 0 logs when disabled")
	}

	logger.Enable()
	logger.Log(result)
	if logger.Size() != 1 {
		t.Error("expected 1 log when enabled")
	}
}

func TestContextLogger_ExportImportJSON(t *testing.T) {
	logger := NewContextLogger()
	logger.Clear()

	result := &BuildResult{
		UsedTokens: 100,
		ContextItems: []*ContextItem{
			{ID: "test", Priority: T0_SystemPrompt, Content: "test content", TokenSize: 10},
		},
	}
	logger.Log(result)

	jsonStr, err := logger.ExportJSON()
	if err != nil {
		t.Fatalf("failed to export JSON: %v", err)
	}

	logger2 := NewContextLogger()
	if err := logger2.ImportJSON(jsonStr); err != nil {
		t.Fatalf("failed to import JSON: %v", err)
	}

	if logger2.Size() != 1 {
		t.Errorf("expected 1 log after import, got %d", logger2.Size())
	}
}

func TestContextLogger_SetEnabled(t *testing.T) {
	logger := NewContextLogger()
	logger.SetEnabled(false)

	result := &BuildResult{
		UsedTokens: 100,
		ContextItems: []*ContextItem{
			{ID: "test", Priority: T0_SystemPrompt, Content: "test", TokenSize: 10},
		},
	}
	logger.Log(result)

	if logger.Size() != 0 {
		t.Error("expected 0 logs when SetEnabled(false)")
	}

	logger.SetEnabled(true)
	logger.Log(result)
	if logger.Size() != 1 {
		t.Error("expected 1 log when SetEnabled(true)")
	}
}

func TestContextItemLog(t *testing.T) {
	item := &ContextItem{
		ID:        "test-id",
		Priority:  T1_RetrievedKnowledge,
		Content:   "Test content",
		TokenSize: 42,
		Metadata:  map[string]any{"key": "value"},
	}

	log := ContextItemLog{
		ID:        item.ID,
		Priority:  item.Priority.String(),
		Content:   item.Content,
		TokenSize: item.TokenSize,
		Metadata:  item.Metadata,
	}

	if log.ID != "test-id" {
		t.Errorf("expected ID 'test-id', got '%s'", log.ID)
	}
	if log.Priority != "T1_RetrievedKnowledge" {
		t.Errorf("expected Priority 'T1_RetrievedKnowledge', got '%s'", log.Priority)
	}
	if log.TokenSize != 42 {
		t.Errorf("expected TokenSize 42, got %d", log.TokenSize)
	}
}