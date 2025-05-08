package types

import (
	"encoding/json"
	"testing"
	"time"
)

type User struct {
	Extends
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) NewMap() (map[string]any, error) {
	return u.Extends.NewMap(u)
}

func (u *User) MarshalJSON() ([]byte, error) {
	return u.Extends.MarshalJSON(u)
}

func TestUser_SetValue(t *testing.T) {
	user := NewUser()
	user.Name = "1"
	user.Age = 100
	user.ESet("d2", "address")
	user.ESet("d1", "address")
	user.ESet("time", time.Now())
	if bs, err := json.Marshal(user); err == nil {
		t.Log(string(bs))
	} else {
		t.Error(err)
	}
}
