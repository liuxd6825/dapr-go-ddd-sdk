package context

import (
	"context"
	"sync"
)

type ContextManager struct {
	config    *Config
	budget    *TokenBudget
	priority  *PriorityManager
	truncator *Truncator
	mu        sync.RWMutex
}

func NewContextManager(config *Config) (*ContextManager, error) {
	budget, err := NewTokenBudget(config.MaxTotalTokens, config.ReservedTokens)
	if err != nil {
		return nil, err
	}

	truncator := NewTruncator(budget, config.MaxHistoryTurns)

	return &ContextManager{
		config:    config,
		budget:    budget,
		priority:  NewPriorityManager(),
		truncator: truncator,
	}, nil
}

func (cm *ContextManager) Build(ctx context.Context, param *BuildParam) (*BuildResult, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.priority.Clear()
	cm.budget.Reset()

	cm.addSystemPrompt()

	history := cm.processHistory(param.ConversationHistory)

	cm.addRetrievedContext(param.RetrievedContext)

	cm.addUserQuery(param.Query)

	items := cm.priority.SortForAssembly()
	templateData := cm.buildTemplateData(param, history)
	finalPrompt, err := RenderPrompt(templateData)
	if err != nil {
		return nil, err
	}

	return &BuildResult{
		FinalPrompt:    finalPrompt,
		UsedTokens:     cm.budget.CountTokens(finalPrompt),
		ContextItems:   items,
		QueryRewritten: param.Query,
	}, nil
}

func (cm *ContextManager) addSystemPrompt() {
	if cm.config.SystemPrompt == "" {
		return
	}

	item := &ContextItem{
		Priority:  T0_SystemPrompt,
		Content:   cm.config.SystemPrompt,
		TokenSize: cm.budget.CountTokens(cm.config.SystemPrompt),
		ID:        "system-prompt",
		Metadata:  map[string]any{"type": "system"},
	}
	cm.priority.Add(item)
	cm.budget.Consume(item.TokenSize)
}

func (cm *ContextManager) processHistory(history []*Content) []*Content {
	if len(history) == 0 {
		return nil
	}

	truncated := cm.truncator.TruncateConversationHistory(history)

	var historyTokens int
	for _, h := range truncated {
		tokens := cm.budget.CountTokens(h.Content)
		historyTokens += tokens
	}

	ragBudget := cm.budget.Allocate(cm.config.RAGRatio)
	historyBudget := cm.budget.Allocate(cm.config.HistoryRatio)

	var selectedHistory []*Content
	var usedTokens int
	for _, h := range truncated {
		hTokens := cm.budget.CountTokens(h.Content)
		if usedTokens+hTokens > historyBudget {
			break
		}
		selectedHistory = append(selectedHistory, h)
		usedTokens += hTokens
	}

	for i, h := range selectedHistory {
		item := &ContextItem{
			Priority:  T2_RecentHistory,
			Content:   h.Content,
			TokenSize: cm.budget.CountTokens(h.Content),
			ID:        "history-" + h.Role + "-" + string(rune(i)),
			Metadata: map[string]any{
				"type":     "history",
				"role":     h.Role,
				"tokens":   cm.budget.CountTokens(h.Content),
				"rag_budget_remaining": ragBudget - usedTokens,
			},
		}
		cm.priority.Add(item)
		cm.budget.Consume(item.TokenSize)
	}

	return selectedHistory
}

func (cm *ContextManager) addRetrievedContext(chunks []string) {
	if len(chunks) == 0 {
		return
	}

	ragBudget := cm.budget.Allocate(cm.config.RAGRatio)
	var usedTokens int

	for i, chunk := range chunks {
		chunkTokens := cm.budget.CountTokens(chunk)
		if usedTokens+chunkTokens > ragBudget {
			break
		}

		item := &ContextItem{
			Priority:  T1_RetrievedKnowledge,
			Content:   chunk,
			TokenSize: chunkTokens,
			ID:        "rag-chunk-" + itoa(i),
			Metadata: map[string]any{
				"type":     "rag",
				"index":    i,
				"total":    len(chunks),
			},
		}
		cm.priority.Add(item)
		cm.budget.Consume(item.TokenSize)
		usedTokens += chunkTokens
	}
}

func (cm *ContextManager) addUserQuery(query string) {
	item := &ContextItem{
		Priority:  T0_UserQuery,
		Content:   query,
		TokenSize: cm.budget.CountTokens(query),
		ID:        "user-query",
		Metadata:  map[string]any{"type": "query"},
	}
	cm.priority.Add(item)
	cm.budget.Consume(item.TokenSize)
}

func (cm *ContextManager) buildTemplateData(param *BuildParam, history []*Content) *TemplateData {
	return &TemplateData{
		SystemPrompt:     cm.config.SystemPrompt,
		Summary:          "",
		History:          history,
		RetrievedContext: param.RetrievedContext,
		CurrentQuery:     param.Query,
	}
}

func (cm *ContextManager) GetConfig() *Config {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config
}

func (cm *ContextManager) GetBudget() *TokenBudget {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.budget
}

func (cm *ContextManager) GetTruncator() *Truncator {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.truncator
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(buf[pos:])
}