package readexcel

type Review struct {
	Columns   []string     `json:"columns"`
	Items     []ReviewItem `json:"items"`
	Sheets    []*Sheet     `json:"sheets"`
	OpenSheet string       `json:"openSheet"`
}

type Sheet struct {
	Name   string `json:"name"`
	MaxCol int    `json:"maxCol"`
	MaxRow int    `json:"maxRow"`
}

type ReviewItem = map[string]any

func (r *Review) AddItems(item ...ReviewItem) {
	r.Items = append(r.Items, item...)
}

func (r *Review) AddColumns(column ...string) {
	r.Columns = append(r.Columns, column...)
}

func (r *Review) AddSheetName(sheet ...*Sheet) {
	r.Sheets = append(r.Sheets, sheet...)
}
