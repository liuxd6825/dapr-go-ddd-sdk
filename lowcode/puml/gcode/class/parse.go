package class

import (
	"bufio"
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/puml/gcode/class/utils"
	"os"
	"path/filepath"
	"strings"
)

type ReadLineOption struct {
	TrimSpace bool
}

type ReadLineOptions = func(opt *ReadLineOption)

type ParseReader interface {
	Scan() bool
	ReadLine(opts ...ReadLineOptions) string
}

type parseReader struct {
	scanner *bufio.Scanner
}

func (p *parseReader) Scan() bool {
	return p.scanner.Scan()
}

func (p *parseReader) ReadLine(opts ...ReadLineOptions) string {
	opt := &ReadLineOption{TrimSpace: true}
	for _, fun := range opts {
		fun(opt)
	}
	txt := p.scanner.Text()
	if opt.TrimSpace {
		txt = strings.Trim(txt, " ")
	}
	if utils.IsHtml(txt) {
		txt, _ = utils.GetHtmlText(txt)
	}
	return utils.StringClear(txt, "\n", "\r")
}

func NewParseReader(filename string) (ParseReader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	return &parseReader{
		scanner: scanner,
	}, nil
}

func ParseFile(ctx context.Context, fileName string) (*UMLBlock, error) {
	filename, err := filepath.Abs(fileName)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	reader, err := NewParseReader(filename)
	if err != nil {
		return nil, err
	}

	block := NewUMLBlock(filename)
	err = block.Parse(ctx, reader)
	return block, err
}
