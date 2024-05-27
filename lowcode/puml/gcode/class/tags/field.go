package tags

type Field struct {
	Name     string `json:"name"`
	DataType string `json:"dataType"`
}

type Fields []*Field

func (f Fields) Add(field ...*Field) Fields {
	return append(f, field...)
}

func (f Fields) Has(name string) bool {
	for _, item := range f {
		if item.Name == name {
			return true
		}
	}
	return false
}

func (f Fields) Find(name string) *Field {
	for _, item := range f {
		if item.Name == name {
			return item
		}
	}
	return nil
}
