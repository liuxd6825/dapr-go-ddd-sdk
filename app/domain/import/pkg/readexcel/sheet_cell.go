package readexcel

import (
	"fmt"
	"time"

	"github.com/tealeg/xlsx"
)

type Cell interface {
	IsTime() bool
	GetTime(date1904 bool) (t time.Time, err error)
	String() string
	Type() xlsx.CellType
	NumFmt() string
	Value() string
}
type sheetCell struct {
	cell *xlsx.Cell
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
