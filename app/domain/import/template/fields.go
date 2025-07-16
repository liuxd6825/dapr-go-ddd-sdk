package template

type TempCreateField struct {
	Id        string           `json:"id"  title:"ID"`
	CaseId    string           `json:"caseId" title:"案件ID"`
	Name      string           `json:"name" title:"名称"`
	BankName  string           `json:"bankName" title:"银行名称"`
	FileId    string           `json:"fileId" title:"文件ID"`
	FileName  string           `json:"fileName"  title:"文件名称"`
	SheetName string           `json:"sheetName"  title:"页名称"`
	SchemaId  string           `json:"schemaId"  title:"架构ID"`
	MapHeads  []*MapHead       `json:"mapHeads"  title:"表头"`
	Fields    []*TemplateField `json:"fields"  title:"字段"`
	Remark    string           `json:"remark" desc:"备注"` // 备注
}

type TempDeleteField struct {
	Id string `json:"id"`
}

type TempUpdateField struct {
	Id        string           `json:"id" title:"ID"`
	FileId    string           `json:"fileId" title:"文件ID"`
	FileName  string           `json:"fileName" title:"文件名称"`
	CaseId    string           `json:"caseId" title:"案件ID"`
	Name      string           `json:"name"  title:"名称"`
	BankName  string           `json:"bankName" title:"银行名称"`
	SheetName string           `json:"sheetName" title:"页名称"`
	SchemaId  string           `json:"schemaId" title:"架构ID"`
	MapHeads  []*MapHead       `json:"mapHeads" title:"表头"`
	Fields    []*TemplateField `json:"fields" title:"字段"`
	Remark    string           `json:"remark"  title:"备注"`
}
