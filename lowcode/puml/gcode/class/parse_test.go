package class

import (
	"context"
	"encoding/json"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/puml/gcode/class/utils"
	"testing"
)

func TestParseFile(t *testing.T) {
	parseFile(t, "./test/example/tag.puml")
}

func TestParseFile_Note(t *testing.T) {
	parseFile(t, "./test/example/note.puml")
}

func TestMakeCode(t *testing.T) {
	ctx := context.Background()
	err := MakeCode(ctx, "./test/example/go.puml", "./template/golang", "./source", "go", option)
	if err != nil {
		t.Error(err)
	}
}

func TestMakeCode_MethodHtml(t *testing.T) {
	ctx := context.Background()
	err := MakeCode(ctx, "./test/example/method_html.puml", "./template/golang", "./source", "go", option)
	if err != nil {
		t.Error(err)
	}
}

func option(opt *MackCodeOption) {
	opt.SetFileName = func(name string) string {
		return utils.SnakeString(name)
	}
}

func parseFile(t *testing.T, umlFile string) {
	ctx := context.Background()
	if black, err := ParseFile(ctx, umlFile); err != nil {
		t.Error(err)
	} else {
		s, err := json.Marshal(black)
		if err != nil {
			t.Error(err)
			return
		}
		println(string(s))
	}
}
