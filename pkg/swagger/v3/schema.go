package swagger3

type ISchema interface {
	GetName() string
	GetType() any
	GetProperties() Properties
	SetProperties(Properties)
	GetRequired() []string
	SetRequired([]string)
	AddRequired(string)
	GetItems() *Property
	Init(schema ISchema)
}

type Schema struct {
	FileName   string     `json:"fileName,omitempty"`
	Ref        string     `json:"$ref,omitempty"`
	Id         string     `json:"$id,omitempty"`
	Schema     string     `json:"$schema,omitempty"`
	Type       any        `json:"type,omitempty"`
	Title      string     `json:"title,omitempty"`
	Name       string     `json:"name,omitempty"`
	Properties Properties `json:"properties,omitempty"`
	Items      *Property  `json:"items,omitempty"`
	Required   []string   `json:"required,omitempty"`
	ReadOnly   bool       `json:"readOnly,omitempty"`
	WriteOnly  bool       `json:"writeOnly,omitempty"`
}

func (s *Schema) GetName() string {
	return s.Name
}

func (s *Schema) GetType() any {
	return s.Type
}

func (s *Schema) GetProperties() Properties {
	return s.Properties
}

func (s *Schema) SetProperties(properties Properties) {
	s.Properties = properties
}

func (s *Schema) GetRequired() []string {
	return s.Required
}

func (s *Schema) SetRequired(strings []string) {
	s.Required = strings
}

func (s *Schema) AddRequired(s2 string) {
	s.Required = append(s.Required, s2)
}

func (s *Schema) GetItems() *Property {
	return s.Items
}

func (s *Schema) Init(schema ISchema) {
	s.Schema = schema.GetName()
}

func newSchema() ISchema {
	return &Schema{}
}
