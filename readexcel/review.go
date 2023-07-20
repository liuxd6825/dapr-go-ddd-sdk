package readexcel

type Review struct {
	Columns []string     `json:"columns"`
	Items   []ReviewItem `json:"items"`
}

type ReviewItem = map[string]any

func (r *Review) AddItems(item ...ReviewItem) {
	r.Items = append(r.Items, item...)
}

func (r *Review) AddColumns(column ...string) {
	r.Columns = append(r.Columns, column...)
}
