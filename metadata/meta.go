package metadata

type GetMetadata interface {
	GetMetadata() Metadata
}

type Metadata interface {
	GetProperties() map[string]Property
	GetProperty(name string) (Property, bool)
	AddProperty(property Property) error
}

type metadata struct {
	properties map[string]Property
}

func NewMetadata(properties []Property) Metadata {
	m := &metadata{}
	for _, item := range properties {
		_ = m.AddProperty(item)
	}
	return m
}

func (m *metadata) GetProperties() map[string]Property {
	return m.properties
}

func (m *metadata) AddProperty(p Property) error {
	m.properties[p.Name()] = p
	return nil
}

func (m *metadata) GetProperty(name string) (Property, bool) {
	p, ok := m.properties[name]
	return p, ok
}
