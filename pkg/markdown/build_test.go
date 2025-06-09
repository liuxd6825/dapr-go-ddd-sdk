package markdown

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/randomutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"testing"
)

func Test_Build(t *testing.T) {
	b := NewBuilder()
	sch := xtest.GetHumanSchema()
	humanName := randomutils.NameCN()
	humanMap := xtest.NewHumanMap(humanName)
	sb, err := b.Build(humanMap, sch)
	if err != nil {
		t.Error(err)
	}
	println(sb.String())
}

func Test_BuildList(t *testing.T) {
	b := NewBuilder()
	sch := xtest.GetHumanSchema()
	var humanMapList []map[string]any
	for i := 0; i < 10; i++ {
		humanName := randomutils.NameCN()
		humanMap := xtest.NewHumanMap(humanName)
		humanMapList = append(humanMapList, humanMap)
	}

	sb, err := b.Build(humanMapList, sch)
	if err != nil {
		t.Error(err)
	}
	println(sb.String())
}
