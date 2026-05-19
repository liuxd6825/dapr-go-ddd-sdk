package context

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

type ContextLog struct {
	Timestamp    time.Time         `json:"timestamp"`
	Level        string            `json:"level"`
	RequestID    string            `json:"requestId"`
	UsedTokens   int               `json:"usedTokens"`
	MaxTokens    int               `json:"maxTokens"`
	TokenUsage   float64           `json:"tokenUsageRatio"`
	ItemsCount   int               `json:"itemsCount"`
	Items        []ContextItemLog  `json:"items"`
	BuildResult  BuildResultLog    `json:"buildResult"`
	QueryRewritten string           `json:"queryRewritten,omitempty"`
}

type ContextItemLog struct {
	ID        string            `json:"id"`
	Priority  string            `json:"priority"`
	Content   string            `json:"content"`
	TokenSize int               `json:"tokenSize"`
	Metadata  map[string]any    `json:"metadata,omitempty"`
}

type BuildResultLog struct {
	FinalPromptLength int      `json:"finalPromptLength"`
	ContextItemsCount int      `json:"contextItemsCount"`
	TokenBudgetUsed   int      `json:"tokenBudgetUsed"`
}

type ContextLogger struct {
	mu         sync.RWMutex
	logs       []*ContextLog
	maxLogs    int
	level      LogLevel
	requestID  string
	enabled    bool
}

func NewContextLogger() *ContextLogger {
	return &ContextLogger{
		logs:     make([]*ContextLog, 0),
		maxLogs:  1000,
		level:    LogLevelInfo,
		enabled:  true,
	}
}

func (cl *ContextLogger) Log(result *BuildResult) {
	if !cl.enabled {
		return
	}

	cl.mu.Lock()
	defer cl.mu.Unlock()

	itemLogs := make([]ContextItemLog, 0, len(result.ContextItems))
	for _, item := range result.ContextItems {
		itemLogs = append(itemLogs, ContextItemLog{
			ID:        item.ID,
			Priority:  item.Priority.String(),
			Content:   item.Content,
			TokenSize: item.TokenSize,
			Metadata:  item.Metadata,
		})
	}

	log := &ContextLog{
		Timestamp:     time.Now(),
		Level:         LogLevelInfo.String(),
		RequestID:     cl.requestID,
		UsedTokens:    result.UsedTokens,
		ItemsCount:    len(result.ContextItems),
		Items:         itemLogs,
		QueryRewritten: result.QueryRewritten,
	}

	cl.logs = append(cl.logs, log)

	if len(cl.logs) > cl.maxLogs {
		cl.logs = cl.logs[1:]
	}
}

func (cl *ContextLogger) LogWithLevel(level LogLevel, result *BuildResult) {
	if level < cl.level {
		return
	}

	cl.mu.Lock()
	defer cl.mu.Unlock()

	itemLogs := make([]ContextItemLog, 0, len(result.ContextItems))
	for _, item := range result.ContextItems {
		itemLogs = append(itemLogs, ContextItemLog{
			ID:        item.ID,
			Priority:  item.Priority.String(),
			Content:   item.Content,
			TokenSize: item.TokenSize,
			Metadata:  item.Metadata,
		})
	}

	log := &ContextLog{
		Timestamp:     time.Now(),
		Level:         level.String(),
		RequestID:     cl.requestID,
		UsedTokens:    result.UsedTokens,
		ItemsCount:    len(result.ContextItems),
		Items:         itemLogs,
		QueryRewritten: result.QueryRewritten,
	}

	cl.logs = append(cl.logs, log)

	if len(cl.logs) > cl.maxLogs {
		cl.logs = cl.logs[1:]
	}
}

func (cl *ContextLogger) GetLogs() []*ContextLog {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	result := make([]*ContextLog, len(cl.logs))
	copy(result, cl.logs)
	return result
}

func (cl *ContextLogger) GetLastLog() *ContextLog {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	if len(cl.logs) == 0 {
		return nil
	}
	return cl.logs[len(cl.logs)-1]
}

func (cl *ContextLogger) Clear() {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.logs = make([]*ContextLog, 0)
}

func (cl *ContextLogger) SetLevel(level LogLevel) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.level = level
}

func (cl *ContextLogger) SetMaxLogs(max int) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.maxLogs = max
}

func (cl *ContextLogger) SetRequestID(id string) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.requestID = id
}

func (cl *ContextLogger) Enable() {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.enabled = true
}

func (cl *ContextLogger) Disable() {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.enabled = false
}

func (cl *ContextLogger) ExportJSON() (string, error) {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	data, err := json.MarshalIndent(cl.logs, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal logs: %w", err)
	}

	return string(data), nil
}

func (cl *ContextLogger) ImportJSON(jsonStr string) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	var logs []*ContextLog
	if err := json.Unmarshal([]byte(jsonStr), &logs); err != nil {
		return fmt.Errorf("failed to unmarshal logs: %w", err)
	}

	cl.logs = logs
	return nil
}

func (cl *ContextLogger) Size() int {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	return len(cl.logs)
}

func (cl *ContextLogger) SetEnabled(enabled bool) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.enabled = enabled
}