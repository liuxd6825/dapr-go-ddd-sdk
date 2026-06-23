package model

import (
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

// Human
// @Description: 人员基本信息
type Human struct {
	xbase.BaseModel `bson:",inline"`

	Code       string     `json:"code" gorm:"code" bson:"code" title:"编码"`
	Name       string     `json:"name" gorm:"name" bson:"name" index:"" title:"姓名"`
	FormerName string     `json:"formerName" gorm:"former_name" bson:"former_name" title:"曾用名"`
	Birthdate  *time.Time `json:"birthday" gorm:"birthdate" bson:"birthdate" title:"出生日期"`
	Age        int        `json:"age" gorm:"age" bson:"age" title:"年龄"`
	PersonType []string   `json:"personType" gorm:"type:json;column:person_type"  bson:"person_type" title:"人员类型"`
	Gender     string     `json:"gender" gorm:"gender" bson:"gender" title:"性别"`
	Education  string     `json:"education" gorm:"education" bson:"education" title:"学历"`
	School     string     `json:"school" gorm:"school" bson:"school" title:"毕业院校"`
	TagId      string     `json:"tagId" gorm:"tag_id" bson:"tag_id" title:"标签ID"`
	TagName    string     `json:"tagName" gorm:"tag_name" bson:"tag_name" title:"标签"`
	TagColor   string     `json:"tagColor" gorm:"tag_color" bson:"tag_color" title:"标签颜色"`
	Remark     string     `json:"remark" gorm:"remark" bson:"remark" title:"备注"`
}

func NewHuman() *Human {
	return &Human{}
}
