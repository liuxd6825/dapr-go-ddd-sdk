package readexcel

type Review struct {
	Columns []string
	Items   []ReviewItem
}

type ReviewItem = map[string]string

func (r *Review) AddItems(item ...ReviewItem) {
	r.Items = append(r.Items, item...)
}

func (r *Review) AddColumns(column ...string) {
	r.Columns = append(r.Columns, column...)
}
