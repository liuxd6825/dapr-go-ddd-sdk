package query

// ExcelSheet

type ExcelSheetFindByNameQuery struct {
	FileId string `json:"fileId"`
	Name   string `json:"name"`
}

type ExcelSheetFindByFieldIdQuery struct {
	FileId string `json:"fileId"`
}

type ExcelSheetFindByDocFieldIdQuery struct {
	DocFileId string `json:"docFileId" query:"doc-file-id" param:"doc-file-id"`
}

type ExcelSheetFindByIdQuery struct {
	Id string `json:"id" query:"id" param:"id" path:"id"`
}
