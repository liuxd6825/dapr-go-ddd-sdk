package doc_extract

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	testing "testing"
)

func TestExtract_docx(t *testing.T) {
	e := newExtract()
	fileName := "./test_file/1.docx"
	if outTxt, err := e.Extract(newFs(fileName), fileName); err != nil {
		t.Error(err)
		return
	} else {
		t.Log(outTxt)
	}
}
func TestExtract_txt(t *testing.T) {
	e := newExtract()
	fileName := "./test_file/2.txt"
	if outTxt, err := e.Extract(newFs(fileName), fileName); err != nil {
		t.Error(err)
		return
	} else {
		t.Log(outTxt)
	}

}

func TestExtract_xlsx(t *testing.T) {
	e := newExtract()
	fileName := "./test_file/3.xlsx"
	if outTxt, err := e.Extract(newFs(fileName), fileName); err != nil {
		t.Error(err)
		return
	} else {
		t.Log(outTxt)
	}

}

func TestExtract_pdf(t *testing.T) {
	e := newExtract()
	fileName := "./test_file/5.pdf"
	if outTxt, err := e.Extract(newFs(fileName), fileName); err != nil {
		t.Error(err)
		return
	} else {
		t.Log(outTxt)
	}
}

func newFs(filename string) afero.Fs {
	fs := afero.NewOsFs()
	fs.Open(filename)
	return fs
}

func newExtract() *Extract {
	logger := logrus.New()
	return NewExtract(logger)
}
