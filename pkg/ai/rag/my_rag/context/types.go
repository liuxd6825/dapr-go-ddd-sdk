package context

import "time"

type Priority int

const (
	T0_SystemPrompt       Priority = 0
	T0_UserQuery          Priority = 1
	T1_RetrievedKnowledge Priority = 2
	T2_RecentHistory      Priority = 3
	T3_LongTermHistory    Priority = 4
)

func (p Priority) IsHead() bool { return p == T0_SystemPrompt }
func (p Priority) IsTail() bool  { return p == T0_UserQuery }

func (p Priority) String() string {
	switch p {
	case T0_SystemPrompt:
		return "T0_SystemPrompt"
	case T1_RetrievedKnowledge:
		return "T1_RetrievedKnowledge"
	case T2_RecentHistory:
		return "T2_RecentHistory"
	case T3_LongTermHistory:
		return "T3_LongTermHistory"
	default:
		if p == T0_UserQuery {
			return "T0_UserQuery"
		}
		return "Unknown"
	}
}

type ContextItem struct {
	Priority  Priority
	Content   string
	TokenSize int
	Metadata  map[string]any
	ID        string
}

func NewContextItem(priority Priority, content, id string) *ContextItem {
	return &ContextItem{
		Priority: priority,
		Content:  content,
		ID:       id,
		Metadata: make(map[string]any),
	}
}

type Config struct {
	MaxTotalTokens    int
	ReservedTokens   int
	SystemPrompt     string
	InitialSummary   string
	RAGRatio         float64
	HistoryRatio     float64
	MaxHistoryTurns  int
	MaxHistoryTokens int
	EnableSummary    bool
	SummaryThreshold int
	SummaryModel     string
	EnableQueryRewrite bool
	EnableReranker   bool
}

func DefaultConfig() *Config {
	return &Config{
		MaxTotalTokens:   8000,
		ReservedTokens:   500,
		SystemPrompt:     "你是一个知识助手，根据以下上下文回答问题：",
		RAGRatio:         0.7,
		HistoryRatio:     0.3,
		MaxHistoryTurns:  5,
		MaxHistoryTokens: 4000,
		EnableSummary:    true,
		SummaryThreshold: 2000,
		EnableQueryRewrite: true,
		EnableReranker:   false,
	}
}

type BuildParam struct {
	Query               string
	RetrievedContext    []string
	ConversationHistory []*Content
}

type BuildResult struct {
	FinalPrompt    string
	UsedTokens     int
	ContextItems   []*ContextItem
	QueryRewritten string
}

type Content struct {
	Role    string
	Content string
	Time    time.Time
}

func (c *Content) GetRole() string {
	return c.Role
}

func (c *Content) GetContent() string {
	return c.Content
}

func NewUserContent(content string) *Content {
	return &Content{
		Role:    "user",
		Content: content,
		Time:    time.Now(),
	}
}

func NewAssistantContent(content string) *Content {
	return &Content{
		Role:    "assistant",
		Content: content,
		Time:    time.Now(),
	}
}

type ExtractedEntities struct {
	UserProfile map[string]string
	KeyFacts    []string
	Preferences map[string]any
}

type SummaryConfig struct {
	Threshold int
	Model     string
}