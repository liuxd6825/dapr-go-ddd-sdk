package server

import (
	"github.com/dop251/goja"
	"strings"
	"testing"
)

type Data struct {
	User User
}
type User struct {
	Name    string
	Age     int
	Email   string
	Address Address
}
type Address struct {
	Street string
	City   string
	State  string
}

func Test_NewObject(t *testing.T) {
	rt := goja.New()
	obj := &Data{
		User: User{
			Name:  "name",
			Age:   10,
			Email: "email",
			Address: Address{
				Street: "street",
				City:   "city",
				State:  "state",
			},
		},
	}
	v, err := NewObject(rt, obj)
	if err != nil {
		t.Fatal(err)
	}
	_ = rt.Set("o1", v)
	_ = rt.Set("o2", obj)
	name, err := rt.RunString(`
		"City="+o1.User.Address.City + " Street="+ o2.User.Address.Street
	`)
	if goja.IsUndefined(name) || goja.IsNull(name) {
		t.Fatal("name is undefined or empty")
		return
	}
	if err != nil {
		t.Fatal(err)
		return
	}
	if name.String() != "City=city Street=street" {
		t.Fatal("name is wrong")
		return
	}
	t.Log(name)
}

func Test_ClasFunc(t *testing.T) {
	fun := func(args ...string) (string, error) {
		return strings.Join(args, "+"), nil
	}
	rt := goja.New()
	_ = rt.Set("join", fun)
	name, err := rt.RunString(`
		join("a","b","c")
	`)
	if goja.IsUndefined(name) || goja.IsNull(name) {
		t.Fatal("name is undefined or empty")
		return
	}
	if err != nil {
		t.Fatal(err)
		return
	}
	if name.String() != "a+b+c" {
		t.Fatal("name is wrong")
		return
	}
	t.Log(name)
}
