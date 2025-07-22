package query

// ExcelRow

type ExcelRowFindQuery struct {
	CaseId    string `json:"caseId"  param:"caseId"  title:"案件Id"`
	DocId     string `json:"docId"  param:"docId" title:"文档Id"`
	FileId    string `json:"fileId"  param:"fileId" title:"文件Id"`
	SheetName string `json:"sheetName"  param:"sheetName" title:"工作页"`
}

type ExcelRowFindByIdQuery struct {
	Id string `json:"id" title:"Id" required:"true"`
}

type ExcelRowFindPreviewQuery struct {
	CaseId  string `json:"caseId" title:"案件Id" required:"true"`
	FileId  string `json:"fileId" title:"文件Id" required:"true"`
	SheetId string `json:"sheetId" title:"SheetId"  required:"true"`
	MaxRows int64  `json:"maxRows" title:"最大记录数"  required:"true"`
}

type ExcelRowFindPreviewQueryResult struct {
	Columns   []string         `json:"columns" title:"列头"`
	SheetName string           `json:"sheetName" title:"工作表"`
	MaxRow    int64            `json:"maxRow" title:"最大行"`
	MaxCol    int64            `json:"maxCol" title:"最大列"`
	Rows      []map[string]any `json:"rows" title:"数据行"`
}
