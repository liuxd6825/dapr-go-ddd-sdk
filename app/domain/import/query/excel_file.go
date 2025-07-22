package query

// ExcelFile

type ExcelFileFindByIdQuery struct {
	FileId string `json:"fileId" param:"fieldId" title:"文件Id" path:"id" `
}
