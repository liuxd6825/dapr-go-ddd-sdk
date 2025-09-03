package field

type HomeCreateField struct {
	Id     string `json:"id"  title:"ID"`
	CaseId string `json:"caseId" title:"案件ID"`
	Name   string `json:"name" title:"名称"`
	Remark string `json:"remark" desc:"备注"` // 备注
}

type HomeDeleteField struct {
	Id string `json:"id"`
}

type HomeUpdateField struct {
	Id string `json:"id" title:"ID"`

	Remark string `json:"remark"  title:"备注"`
}
