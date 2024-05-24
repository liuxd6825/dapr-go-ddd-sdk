package schema

type Manager struct {
	schemaItems map[string]*Schema
}

var _manager = NewManager()

func NewManager() *Manager {
	return &Manager{schemaItems: map[string]*Schema{}}
}

func GetManager() *Manager {
	return _manager
}

func (m *Manager) Add(schema *Schema) {
	m.schemaItems[schema.Id] = schema
}

func (m *Manager) Get(id string) *Schema {
	return m.schemaItems[id]
}
