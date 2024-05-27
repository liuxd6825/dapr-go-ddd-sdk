package class

import (
	"context"
	"strings"
)

type Note struct {
	BaseElement
	Texts []string `json:"texts"`
	Name  string   `json:"name"`
}

func NewNote(line string) *Note {
	text := []string{line}
	return &Note{Texts: text}
}

func (e *Note) AddText(text string) {
	e.Texts = append(e.Texts, text)
}

func (e *Note) GetName() string {
	return e.Name
}

func (e *Note) Parse(ctx context.Context, reader ParseReader) (bool, error) {
	for reader.Scan() {
		line := reader.ReadLine()
		if strings.HasPrefix(line, "  ") {
			continue
		} else if strings.Contains(line, "end") && strings.Contains(line, "note") {
			return true, nil
		} else {
			return false, nil
		}
	}
	return false, nil
}
