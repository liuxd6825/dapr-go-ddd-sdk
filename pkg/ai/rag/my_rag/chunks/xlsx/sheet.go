package xlsx

type Sheet struct {
	headRow   int
	FileName  string
	SheetName string
	Fields    []string
	Rows      []string
}

func NewSheet() *Sheet {
	return &Sheet{
		FileName:  "",
		SheetName: "",
		headRow:   0,
		Fields:    []string{},
		Rows:      []string{},
	}
}

func (s *Sheet) Init(name string, headRow int, texts []string) {
	s.SheetName = name
	s.headRow = headRow
	if headRow < 0 {
		s.Rows = texts
	} else {
		s.Fields = texts[0 : headRow+1]
		s.Rows = texts[headRow+1 : len(texts)]
	}
}

func (s *Sheet) GetFields() []string {
	return s.Fields
}

func (s *Sheet) GetRows() []string {
	return s.Rows
}
