package context

type Truncator struct {
	budget   *TokenBudget
	maxTurns int
}

func NewTruncator(budget *TokenBudget, maxTurns int) *Truncator {
	return &Truncator{
		budget:   budget,
		maxTurns: maxTurns,
	}
}

func (t *Truncator) TruncateConversationHistory(history []*Content) []*Content {
	if len(history) == 0 {
		return history
	}

	var result []*Content
	var currentTurn []*Content
	remaining := t.budget.Remaining()

	for i := len(history) - 1; i >= 0; i-- {
		item := history[i]
		tokens := t.budget.CountTokens(item.Content)

		if remaining < tokens && len(currentTurn) > 0 {
			break
		}

		currentTurn = append([]*Content{item}, currentTurn...)

		if item.Role == "user" {
			if len(result)+len(currentTurn)/2 > t.maxTurns && len(currentTurn) >= 2 {
				break
			}

			if len(result) >= t.maxTurns && item.Role == "user" {
				break
			}

			result = append(currentTurn, result...)
			currentTurn = nil
		}

		remaining -= tokens
	}

	if len(currentTurn) > 0 {
		result = append(currentTurn, result...)
	}

	return result
}

func (t *Truncator) TruncateByTokens(items []*ContextItem, maxTokens int) []*ContextItem {
	if len(items) == 0 {
		return items
	}

	var result []*ContextItem
	var usedTokens int

	for _, item := range items {
		if item.TokenSize == 0 {
			item.TokenSize = t.budget.CountTokens(item.Content)
		}

		if usedTokens+item.TokenSize > maxTokens {
			continue
		}

		result = append(result, item)
		usedTokens += item.TokenSize
	}

	return result
}

func (t *Truncator) GetMaxTurns() int {
	return t.maxTurns
}

func (t *Truncator) SetMaxTurns(maxTurns int) {
	t.maxTurns = maxTurns
}

func (t *Truncator) GetBudget() *TokenBudget {
	return t.budget
}