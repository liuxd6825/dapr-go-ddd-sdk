package builder

type Field struct {
	Name     string `json:"name"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}
