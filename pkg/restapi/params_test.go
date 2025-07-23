package restapi

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Command struct {
	Id   string `json:"id" validate:"required" title:"主键"`
	Data struct {
		Name string `json:"id" validate:"required" title:"姓名"`
	} `json:"data" validate:"required" title:"数据"`
}

func Test_Validate(t *testing.T) {
	err := Validate(Command{})
	assert.Error(t, err)
	t.Log(err)
}
