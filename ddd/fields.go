package ddd

type Fields struct {
	CreatedTime string
	CreatorId   string
	CreatorName string
	UpdatedTime string
	UpdaterId   string
	UpdaterName string
	DeletedTime string
	DeleterId   string
	DeleterName string
	IsDeleted   string
	TenantId    string
	Id          string
	CaseId      string
}

var fields = newFields()

func GetFields() *Fields {
	return fields
}

func newFields() *Fields {
	return &Fields{
		CreatedTime: "createdTime",
		CreatorId:   "creatorId",
		CreatorName: "creatorName",
		UpdatedTime: "updatedTime",
		UpdaterId:   "updaterId",
		UpdaterName: "updaterName",
		DeletedTime: "deletedTime",
		DeleterId:   "deleterId",
		DeleterName: "deleterName",
		IsDeleted:   "isDeleted",
		TenantId:    "tenantId",
		Id:          "id",
		CaseId:      "caseId",
	}
}
