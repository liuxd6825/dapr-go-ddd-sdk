package query

// ExcelFile

type ExcelFileFindByIdQuery struct {
	FileId string `json:"fileId" param:"file-id" title:"文件Id" path:"id" `
}

type ExcelFileFindByDocFileIdQuery struct {
	DocFileId string `json:"docFileId" param:"doc-file-id" title:"文件Id" query:"doc-file-id" `
}
