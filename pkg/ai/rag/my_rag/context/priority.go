package context

type PriorityManager struct {
	items map[Priority][]*ContextItem
}

func NewPriorityManager() *PriorityManager {
	return &PriorityManager{
		items: make(map[Priority][]*ContextItem),
	}
}

func (pm *PriorityManager) Add(item *ContextItem) {
	pm.items[item.Priority] = append(pm.items[item.Priority], item)
}

func (pm *PriorityManager) AddContent(priority Priority, content, id string) {
	pm.Add(NewContextItem(priority, content, id))
}

func (pm *PriorityManager) GetByPriority(p Priority) []*ContextItem {
	return pm.items[p]
}

func (pm *PriorityManager) GetAll() []*ContextItem {
	return pm.SortForAssembly()
}

func (pm *PriorityManager) TotalTokens(budget *TokenBudget) int {
	var total int
	for _, items := range pm.items {
		for _, item := range items {
			if item.TokenSize == 0 {
				item.TokenSize = budget.CountTokens(item.Content)
			}
			total += item.TokenSize
		}
	}
	return total
}

func (pm *PriorityManager) Clear() {
	pm.items = make(map[Priority][]*ContextItem)
}

func (pm *PriorityManager) SortForAssembly() []*ContextItem {
	var result []*ContextItem

	result = append(result, pm.items[T0_SystemPrompt]...)
	result = append(result, pm.items[T3_LongTermHistory]...)
	result = append(result, pm.items[T2_RecentHistory]...)
	result = append(result, pm.items[T1_RetrievedKnowledge]...)
	result = append(result, pm.items[T0_UserQuery]...)

	return result
}

func (pm *PriorityManager) FilterByTokenLimit(budget *TokenBudget, maxTokens int) []*ContextItem {
	var result []*ContextItem
	var usedTokens int

	for _, item := range pm.SortForAssembly() {
		if item.TokenSize == 0 {
			item.TokenSize = budget.CountTokens(item.Content)
		}
		if usedTokens+item.TokenSize > maxTokens {
			continue
		}
		result = append(result, item)
		usedTokens += item.TokenSize
	}

	return result
}

func (pm *PriorityManager) Merge(other *PriorityManager) {
	for priority, items := range other.items {
		pm.items[priority] = append(pm.items[priority], items...)
	}
}