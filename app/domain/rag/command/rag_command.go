package command

type RagCreateCaseCommand struct {
	CommandId string            `json:"commandId"`
	Data      RagCreateCaseData `json:"data"`
}

type RagCreateCaseData struct {
	CaseId string `json:"caseId"`
}

type RagCreateTenantCommand struct {
	CommandId string              `json:"commandId"`
	Data      RagCreateTenantData `json:"data"`
}

type RagCreateTenantData struct {
	TenantId string `json:"tenantId"`
}
