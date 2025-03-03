package db_pkg

type Filter struct {
}

type FilterField struct {
	Field string
}

func NewFilter() *Filter {
	return &Filter{}
}

func (r *Filter) Field(field string) *Filter {

}
