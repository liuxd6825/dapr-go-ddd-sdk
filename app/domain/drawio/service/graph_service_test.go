package service

import (
	_ "embed"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/drawio/pkg/mxgraph"
	"testing"
)

//go:embed test_file/diff1/3_link_update_target_is_empty.drawio
var linkUpdateTargetIsEmptyDrawio string

//go:embed test_file/diff1/3_link_update_target_is_empty.json
var linkUpdateTargetIsEmptyJson string

func Test_Diff3(t *testing.T) {
	drawFile, err := mxgraph.NewDrawioFile(linkUpdateTargetIsEmptyDrawio)
	if err != nil {
		t.Fatal(err)
	}
	element := drawFile.GetElement("kkcGR8n074Od3FpKltgG-5")
	if element == nil {
		t.Fatal("element is nil")
	}

	diff := mxgraph.NewFileDiff(linkUpdateTargetIsEmptyJson)
	if diff == nil {
		t.Fatal("diff is nil")
	}

	graphService := NewGraphService()
	saveBatch := graphService.GetSaveNodes(drawFile, diff)
	t.Log(saveBatch)
}
