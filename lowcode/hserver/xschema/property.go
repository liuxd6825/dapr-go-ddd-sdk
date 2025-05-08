package xschema

type Properties map[string]*Property

type Property struct {
	Id          string     `json:"$id,omitempty"`
	Name        string     `json:"-"`
	Type        any        `json:"type,omitempty"`
	Title       string     `json:"title,omitempty"`
	Format      Format     `json:"format,omitempty"`
	Pattern     string     `json:"pattern,omitempty"`
	Ref         string     `json:"$ref,omitempty"`
	Description string     `json:"description,omitempty"`
	Properties  Properties `json:"properties,omitempty"`
	Required    []string   `json:"required,omitempty"`
	Items       *Property  `json:"items,omitempty"` // nil or []*Schema or *Schema
	AllOf       []*Schema  `json:"allOf,omitempty"`
	Minimum     *int       `json:"minimum,omitempty"`
	Maximum     *int       `json:"maximum,omitempty"`
	MinLength   *int       `json:"minLength,omitempty"`
	MaxLength   *int       `json:"maxLength,omitempty"`
	ReadOnly    bool       `json:"readOnly,omitempty"`
	WriteOnly   bool       `json:"writeOnly,omitempty"`
	Examples    []any      `json:"examples,omitempty"`
	Deprecated  bool       `json:"deprecated,omitempty"`
	types       []string
}
