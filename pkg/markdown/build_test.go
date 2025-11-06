package markdown

import (
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/randomutils"
	xtest2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
)

func Test_Build(t *testing.T) {
	b := NewBuilder()
	sch := xtest2.GetHumanSchema()
	humanName := randomutils.NameCN()
	humanMap := xtest2.NewHumanMap(humanName)
	sb, err := b.Build(humanMap, sch)
	if err != nil {
		t.Error(err)
	}
	println(sb.String())
}

func Test_BuildList(t *testing.T) {
	b := NewBuilder()
	sch := xtest2.GetHumanSchema()
	var humanMapList []map[string]any
	for i := 0; i < 10; i++ {
		humanName := randomutils.NameCN()
		humanMap := xtest2.NewHumanMap(humanName)
		humanMapList = append(humanMapList, humanMap)
	}

	sb, err := b.Build(humanMapList, sch)
	if err != nil {
		t.Error(err)
	}
	println(sb.String())
}
