package dbschema

import (
	"testing"
	"time"
)

type Human struct {
	Id         string     `gorm:"primaryKey"`
	Name       string     `gorm:"name"`
	Age        int        `gorm:"age"`
	Analyse    string     `gorm:"analyse"`
	Birthday   *time.Time `gorm:"birthday"`
	PeopleType []string   `gorm:"people_type;type:text[]"`
	Tags       []string   `gorm:"tag;type:text[]"`
}

func Test_NewSchemaWithStruct(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Error(err)
		}
	}()
	humanSch := NewSchemaWithStruct("human", &Human{}, "human")
	println(humanSch)
}
