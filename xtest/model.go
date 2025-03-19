package xtest

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/randomutils"
	"time"
)

type Human struct {
	Id         string     `gorm:"primaryKey"`
	Name       string     `gorm:"name"`
	Age        int        `gorm:"age"`
	Analyse    string     `gorm:"analyse"`
	Birthday   *time.Time `gorm:"birthday"`
	PeopleType []string   `gorm:"type:text;serializer:json"`
}

func NewHumanList(count int64, humanName string) []*Human {
	list := make([]*Human, count)
	for i := int64(0); i < count; i++ {
		entity := &Human{
			Id:         randomutils.NewId(),
			Name:       humanName,
			Analyse:    "",
			Age:        randomutils.IntMax(100),
			Birthday:   randomutils.PDate(),
			PeopleType: []string{randomutils.String(10)},
		}
		list[i] = entity
	}
	return list
}
