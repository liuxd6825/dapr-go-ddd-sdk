package xtest

import (
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/randomutils"
)

type Human struct {
	Id          string     `gorm:"primaryKey" json:"id"`
	TenantId    string     `gorm:"tenant_id" json:"tenantId"`
	Name        string     `gorm:"name" json:"name"`
	Age         int        `gorm:"age" json:"age"`
	Analyse     string     `gorm:"analyse" json:"analyse"`
	Birthday    *time.Time `gorm:"birthday,index:,sort:desc," json:"birthday"`
	PeopleType  []string   `gorm:"type:json;,serializer:json" json:"peopleType"`
	CreatedTime time.Time  `gorm:"index:,sort:desc,created_time" json:"createdTime"` // 发生事件
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
		list[i] = NewHumanMap(humanName)
	}
	return list
}

func NewHumanMap(humanName string) map[string]any {
	gender := "男"
	vMax := randomutils.IntMax(10)
	if vMax%2 == 1 {
		gender = "女"
	}
	entity := map[string]any{
		"id":         randomutils.NewId(),
		"tenantId":   "test",
		"name":       humanName,
		"analyse":    "",
		"age":        randomutils.IntMax(100),
		"birthday":   randomutils.PDate(),
		"peopleType": []string{randomutils.String(10)},
		"gender":     gender,
	}
	return entity
}
