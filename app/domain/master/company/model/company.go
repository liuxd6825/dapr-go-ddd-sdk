package model

import (
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

// Company
// @Description: 公司基本信息
type Company struct {
	xbase.BaseModel `bson:",inline"`

	Code             string     `json:"code" gorm:"code" bson:"code" title:"公司编码"`
	Name             string     `json:"name" gorm:"name" bson:"name" index:"" validate:"-" title:"公司名称"`
	RegNo            string     `json:"regNo" gorm:"reg_no" bson:"reg_no" title:"工商注册号"`
	OperatingStatus  string     `json:"operatingStatus" gorm:"operating_status" bson:"operating_status" title:"经营状态"`
	CreditCode       string     `json:"creditCode" gorm:"credit_code" bson:"credit_code" title:"统一社会信用代码"`
	Identification   string     `json:"identification" gorm:"identification" bson:"identification" title:"纳税人识别号"`
	OrgCode          string     `json:"orgCode" gorm:"org_code" bson:"org_code" title:"组织机构代码"`
	BusinessTerm     string     `json:"businessTerm" gorm:"business_term" bson:"business_term" title:"营业期限"`
	Qualification    string     `json:"qualification" gorm:"qualification" bson:"qualification" title:"纳税人资质"`
	ApprovalDate     *time.Time `json:"approvalDate" gorm:"approval_date" bson:"approval_date" title:"核准日期"`
	CreateDate       *time.Time `json:"createDate" gorm:"create_date" bson:"create_date" title:"成立时间"`
	EntType          string     `json:"entType" gorm:"ent_type" bson:"ent_type" title:"企业类型"`
	Industry         string     `json:"industry" gorm:"industry" bson:"industry" title:"行业"`
	PersonnelSize    int        `json:"personnelSize" gorm:"personnel_size" bson:"personnel_size" title:"人员规模"`
	InsuredPersons   string     `json:"insuredPersons" gorm:"insured_persons" bson:"insured_persons" title:"参保人数"`
	RegAuthority     string     `json:"regAuthority" gorm:"reg_authority" bson:"reg_authority" title:"登记机关"`
	NameUsed         string     `json:"nameUsed" gorm:"name_used" bson:"name_used" title:"曾用名"`
	EngName          string     `json:"engName" gorm:"eng_name" bson:"eng_name" title:"英文名称"`
	Addr             string     `json:"addr" gorm:"addr" bson:"addr" title:"注册地址"`
	AddrType         string     `json:"addrType" gorm:"addr_type" bson:"addr_type" title:"地址类型"`
	RegCapitalType   string     `json:"regCapitalType" gorm:"reg_capital_type" bson:"reg_capital_type" title:"注册资本类型"`
	RegCapital       string     `json:"regCapital" gorm:"reg_capital" bson:"reg_capital" title:"注册资本"`
	PaidInCapital    string     `json:"paidInCapital" gorm:"paid_in_capital" bson:"paid_in_capital" title:"实缴资本"`
	NatureOfBusiness string     `json:"natureOfBusiness" gorm:"nature_of_business" bson:"nature_of_business" title:"经营范围"`
	TagId            string     `json:"tagId" gorm:"tag_id" bson:"tag_id" title:"tagId"`
	TagColor         string     `json:"tagColor" gorm:"tag_color" bson:"tag_color" title:"tagColor"`
	TagName          string     `json:"tagName" gorm:"tag_name" bson:"tag_name" title:"标签"`
}

func NewCompany() *Company {
	return &Company{}
}