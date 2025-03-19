package xtest

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/randomutils"
	"time"
)

type Human struct {
	Id         string     `gorm:"primaryKey" json:"id"`
	TenantId   string     `gorm:"tenant_id" json:"tenantId"`
	Name       string     `gorm:"name" json:"name"`
	Age        int        `gorm:"age" json:"age"`
	Analyse    string     `gorm:"analyse" json:"analyse"`
	Birthday   *time.Time `gorm:"birthday" json:"birthday"`
	PeopleType []string   `gorm:"type:text;serializer:json" json:"peopleType"`
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

func NewHumanMapList(count int64, humanName string) []map[string]any {
	list := make([]map[string]any, count)
	for i := int64(0); i < count; i++ {
		entity := map[string]any{
			"id":         randomutils.NewId(),
			"name":       humanName,
			"analyse":    "",
			"age":        randomutils.IntMax(100),
			"birthday":   randomutils.PDate(),
			"peopleType": []string{randomutils.String(10)},
		}
		list[i] = entity
	}
	return list
}
