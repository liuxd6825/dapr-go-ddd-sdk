package readexcel

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"time"
)

type Rows interface {
	Count() int64
	Row(row int) []Cell
}

type Cell interface {
	IsTime() bool
	GetTime(date1904 bool) (t time.Time, err error)
	String() string
	Type() xlsx.CellType
	NumFmt() string
	Value() string
}

type sheetRow struct {
	sheet *xlsx.Sheet
}

type sheetCell struct {
	cell *xlsx.Cell
}

type mapRow struct {
	list []map[string]any
}

type mapCell struct {
	value    any
	isTime   bool
	cellType xlsx.CellType
	numFmt   string
}

func (c *mapCell) IsTime() bool {
	return c.isTime
}

func (c *mapCell) GetTime(date1904 bool) (t time.Time, err error) {
	return time.Now(), nil
}

func (c *mapCell) String() string {
	return fmt.Sprintf("%v", c.value)
}

func (c *mapCell) Type() xlsx.CellType {
	return c.cellType
}

func (c *mapCell) NumFmt() string {
	return c.numFmt
}

func (c *mapCell) Value() string {
	return fmt.Sprintf("%v", c.value)
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

func NewCell(cell *xlsx.Cell) Cell {
	return &sheetCell{
		cell: cell,
	}
}

func NewCellByMap(cellValue any) Cell {
	return &mapCell{
		value: cellValue,
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

func (s *sheetCell) IsTime() bool {
	return s.cell.IsTime()
}

func (s *sheetCell) GetTime(date1904 bool) (t time.Time, err error) {
	return s.cell.GetTime(date1904)
}

func (s *sheetCell) String() string {
	return s.cell.String()
}

func (s *sheetCell) Type() xlsx.CellType {
	return s.cell.Type()
}

func (s *sheetCell) Value() string {
	return s.cell.Value
}

func (s *sheetCell) NumFmt() string {
	return s.cell.NumFmt
}
