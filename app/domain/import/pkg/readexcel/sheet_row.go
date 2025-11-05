package readexcel

import (
	"github.com/tealeg/xlsx"
)

type Rows interface {
	Count() int64
	Row(row int) []Cell
}

type sheetRow struct {
	sheet *xlsx.Sheet
}

type mapRow struct {
	list []map[string]any
}

func NewRowBySheet(sheet *xlsx.Sheet) Rows {
	return &sheetRow{
		sheet: sheet,
	}
}

func NewRowByMap(list []map[string]any) Rows {
	return &mapRow{
		list: list,
	}
}

func (s *sheetRow) Count() int64 {
	return int64(s.sheet.MaxRow)
}

func (s *sheetRow) Row(i int) []Cell {
	var cells []Cell
	for _, cell := range s.sheet.Row(i).Cells {
		cells = append(cells, NewCell(cell))
	}
	return cells
}

func (m *mapRow) Count() int64 {
	return int64(len(m.list))
}

func (m *mapRow) Row(row int) []Cell {
	var cells []Cell
	data := m.list[row]
	length := len(data)
	for i := 0; i < length-1; i++ {
		key := GetCellLabel(i + 1)
		value := data[key]
		cells = append(cells, NewCellByMap(value))
	}
	return cells
}
