package restapp

type Table struct {
	DbKey           string
	TableName       string
	Object          any
	IsMongo         bool
	IsGORM          bool
	IsElasticSearch bool
}

type Tables struct {
	items []*Table
}

func NewTables() *Tables {
	return &Tables{items: []*Table{}}
}

func NewTable(dbKey string, tableName string, object any) *Table {
	return &Table{DbKey: dbKey, TableName: tableName, Object: object}
}

func (t *Table) SetIsMongo(v bool) *Table {
	t.IsMongo = v
	return t
}

func (t *Table) SetIsGORM(v bool) *Table {
	t.IsGORM = v
	return t
}

func (t *Table) SetIsElasticSearch(v bool) *Table {
	t.IsElasticSearch = v
	return t
}

func (s *Tables) Append(dbKey string, tableName string, object any) *Table {
	table := NewTable(dbKey, tableName, object)
	s.items = append(s.items, table)
	return table
}
