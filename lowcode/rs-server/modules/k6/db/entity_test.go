package db

import (
	"fmt"
	"github.com/dop251/goja"
	"testing"
)

func TestEntity(t *testing.T) {
	vm := goja.New()
	vm.Set("setEntity", setEntity)
	vm.RunString("entity=setEntity({id: '123', tenantId: 'test'})")
	e := vm.Get("entity").Export()
	fmt.Println(e)
	if entity, ok := e.(Entity); ok {
		if entity.GetId() != "123" || entity.GetTenantId() != "test" {
			t.Error("entity not set correctly")
		} else {
			fmt.Println("entity set correctly")
		}
	} else {
		t.Error("entity not set correctly")
	}

	data := map[string]interface{}{"id": "123", "tenantId": "test"}
	ent := Entity(data)
	t.Log(ent.GetId())

}

func setEntity(e Entity) Entity {
	fmt.Println(e)
	return e
}
