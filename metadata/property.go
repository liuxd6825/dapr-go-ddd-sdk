package metadata

type Property interface {
	Name() string
	DataType() DataType
}

type property struct {
	name        string
	dataType    DataType
	description string
}

func NewProperty(name string, dataType DataType, description string) Property {
	return &property{
		name:        name,
		dataType:    dataType,
		description: description,
	}
}

func (p *property) Name() string {
	return p.name
}

func (p *property) DataType() DataType {
	return p.dataType
}

func (p *property) Description() string {
	return p.description
}
