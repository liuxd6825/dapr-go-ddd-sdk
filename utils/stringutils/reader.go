package stringutils

import "strings"

type Reader struct {
	lines []string
}

func NewReader(txt string) *Reader {
	lines := strings.Split(txt, "\n")
	return &Reader{lines: lines}
}

func (r *Reader) Lines() []string {
	return r.lines
}
