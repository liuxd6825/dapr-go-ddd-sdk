package mxgraph

import (
	_ "embed"
	"testing"
)

//go:embed test_file/333.drawio
var content string

func Test_NewDrawioFile(t *testing.T) {
	drawFile, err := NewDrawioFile(content)
	if err != nil {
		t.Fatal(err)
	}
	element := drawFile.GetElement("Z4WWRIDEBgNPTWCXEiSJ-1")
	if element == nil {
		t.Fatal("element is nil")
	}
}

//go:embed test_file/diff1/3_link_update_target_is_empty.drawio
var linkUpdateTargetIsEmptyDrawio string

//go:embed test_file/diff1/3_link_update_target_is_empty.json
var linkUpdateTargetIsEmptyJson string

func Test_Diff3(t *testing.T) {
	drawFile, err := NewDrawioFile(linkUpdateTargetIsEmptyDrawio)
	if err != nil {
		t.Fatal(err)
	}
	element := drawFile.GetElement("kkcGR8n074Od3FpKltgG-5")
	if element == nil {
		t.Fatal("element is nil")
	}

	diff := NewFileDiff(linkUpdateTargetIsEmptyJson)
	if diff == nil {
		t.Fatal("diff is nil")
	}
	for _, updatePage := range diff.U {
		for id, u := range updatePage.Cells.U {
			if id != "" {
				el := drawFile.GetElement(id)
				if el != nil {
					t.Log(el, u)
				}
				xmlName := el.GetXMLName()
				switch xmlName.Local {
				case "mxCell":
					{
						mxCell := el.(*MxCell)
						t.Log(mxCell)
						break
					}
				case "UserObject":
					{
						break
					}
				case "object":
					{
						break
					}
				}
			}
		}
	}
}
