package runtime

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
)

type PkgValue struct {
	hostVm  *goja.Runtime
	objIns  *goja.Object
	methods *types.CMap[*goja.Callable]
}

func NewPkgValue(hostVm *goja.Runtime, value goja.Value) *PkgValue {
	objInst := value.ToObject(hostVm)
	return &PkgValue{hostVm: hostVm, objIns: objInst}
}

func getMethods(obj *goja.Object) *types.CMap[*goja.Callable] {
	methods := types.NewCMap[*goja.Callable]()
	Prototype := obj.Prototype()
	protoKeys := Prototype.Keys()
	for _, key := range protoKeys {
		fmt.Println("-", key)
		val := obj.Get(key)
		fun, ok := goja.AssertFunction(val)
		if ok {
			methods.Set(key, &fun)
		}
	}
	return methods
}

func (v *PkgValue) Get(name string) goja.Value {
	return v.objIns.Get(name)
}

func (v *PkgValue) Set(name string, val any) error {
	return v.objIns.Set(name, val)
}

func (v *PkgValue) Has(name string) bool {
	val := v.objIns.Get(name)
	if val == nil {
		return false
	}
	return true
}

func (v *PkgValue) Keys() []string {
	keys := v.objIns.Keys()
	//keys = append(keys, v.objIns.Prototype().Keys()...)
	return keys
}
