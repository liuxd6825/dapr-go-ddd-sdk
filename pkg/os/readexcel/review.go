package readexcel

type Review struct {
	Sheets    []*Sheet `json:"sheets"`
	OpenSheet string   `json:"openSheet"`
}

type Sheet struct {
	Columns []string     `json:"columns"`
	Name    string       `json:"name"`
	MaxCol  int64        `json:"maxCol"`
	MaxRow  int64        `json:"maxRow"`
	Items   []ReviewItem `json:"items"`
}

type ReviewItem = map[string]any

func (r *Sheet) AddItems(item ...ReviewItem) {
	r.Items = append(r.Items, item...)
}

func (r *Sheet) AddColumns(column ...string) {
	r.Columns = append(r.Columns, column...)
}

func (r *Review) AddSheet(sheet ...*Sheet) {
	r.Sheets = append(r.Sheets, sheet...)
}
