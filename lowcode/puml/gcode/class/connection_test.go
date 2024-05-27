package class

import (
	"encoding/json"
	"testing"
)

func TestSetNodes(t *testing.T) {
	c := Connection{}
	c.SetNodes(`AbstractList "*"<|--"1" ArrayList : 驱动 >`)
	bs, _ := json.Marshal(c)
	t.Log(string(bs))
}
