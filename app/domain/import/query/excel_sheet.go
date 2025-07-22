package query

// ExcelSheet

type ExcelSheetFindByNameQuery struct {
	FileId string `json:"fileId"`
	Name   string `json:"name"`
}

type ExcelSheetFindByFieldIdQuery struct {
	FileId string `json:"fileId"`
}

type ExcelSheetFindByIdQuery struct {
	Id string `json:"id" `
}
