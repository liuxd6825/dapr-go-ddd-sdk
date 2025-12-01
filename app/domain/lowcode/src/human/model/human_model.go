package model

import "time"

// BaseModel 模拟 base.json 中的基础字段
// 通常包含 ID, CreateTime, UpdateTime 等
type BaseModel struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

// Human 人员基本信息
// 对应数据库表: human
type Human struct {
	BaseModel

	// CaseId 项目Id
	// Meta: dbField.readable=true, column.hide=true, query.allowQuery=false
	CaseId string `json:"caseId" meta:"dbField.readable=true, column.hide=true, query.allowQuery=false" `

	// Name 姓名 (Required)
	Name string `json:"name" meta:"ui5-input, width=150" title:"姓名" required:"true"`

	// FormerName 曾姓名
	// Type: string or null
	FormerName string `json:"formerName,omitempty" title:"曾姓名" `

	// Birthday 出生日期
	// Type: date or null. 使用 string 存储 YYYY-MM-DD 格式，或者使用 time.Time
	Birthday *time.Time `json:"birthday,omitempty"`

	// Age 年龄
	// Type: integer or null
	// Meta: readonly, notField (通常表示数据库计算字段或非持久化字段)
	Age *int `json:"age,omitempty" gorm:"-"`

	// PersonType 人员类型
	// Type: array of strings
	// Meta: ui5-multi-combobox
	PersonType []string `json:"personType,omitempty" gorm:"type:json;serializer:json"`

	// Gender 性别
	// Type: string or null (Options: man, women)
	Gender *string `json:"gender,omitempty"`

	// Education 学历
	// Type: string or null
	Education *string `json:"education,omitempty"`

	// School 毕业院校
	// Type: string or null
	School *string `json:"school,omitempty"`

	// TagId tagId
	// Meta: hide=true
	TagId *string `json:"tagId,omitempty" gorm:"size:500"`

	// TagColor tagColor
	// Meta: hide=true
	TagColor *string `json:"tagColor,omitempty" gorm:"size:500"`

	// TagName 标签
	// Meta: ui5e-tag-selector
	TagName string `json:"tagName,omitempty" gorm:"size:500"`

	// Remark 备注
	// Meta: colSpan="S1 M2 L2 XL2"
	Remark string `json:"remark,omitempty" gorm:"size:500"`
}

// TableName 指定数据库表名
func (Human) TableName() string {
	return "human"
}
