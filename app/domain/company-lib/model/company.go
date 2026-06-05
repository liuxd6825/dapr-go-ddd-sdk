package model

type Company struct {
	Id            string    `json:"id" `
	Name          string    `json:"name"`          // 名称
	RegNo         string    `json:"reg_no"`        // 工商注册号
	OperStatus    string    `json:"oper_status"`   // 经营状态
	CreditCode    string    `json:"credit_code"`   // 统一社会信用代码
	Iden          string    `json:"iden"`          // 纳税人识别号
	ApprDate      string    `json:"appr_date"`     // 核准日期
	CreateDate    string    `json:"create_date"`   // 成立时间
	Qual          string    `json:"qual"`          // 纳税人资质
	EntType       string    `json:"ent_type"`      // 企业类型
	RegAuthority  string    `json:"reg_authority"` // 登记机关
	EngName       string    `json:"eng_name"`      // 英文名称
	Addr          string    `json:"addr"`          // 注册地址
	SpecificAddr  string    `json:"specific_addr"` // 详细地址
	LegalPerson   string    `json:"legal_person"`  // 法人
	BusInfo       RawJSON `json:"bus_info"`      // 工商信息
	ShInfo        RawJSON `json:"sh_info"`       // 股东信息
	KeyPerson     RawJSON `json:"key_person"`    // 主要人员
	IsTransform   bool      `json:"is_transform"`
	SourceUrl     string    `json:"source_url"`
	SourceId      string    `json:"source_id"`
	SourceSystem  string    `json:"source_system"`
	Remark        string    `json:"remark"`
	CreatedAt     ESTime    `json:"created_at"`
	UpdatedAt     ESTime    `json:"updated_at"`
	PrevUpdatedAt ESTime    `json:"prev_updated_at"`
}
