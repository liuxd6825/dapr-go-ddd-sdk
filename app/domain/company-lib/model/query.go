package model

type CompanyQuery struct {
	Name         string `json:"name" param:"name" query:"name"`
	RegNo        string `json:"regNo" param:"regNo" query:"reg-no"`
	OperStatus   string `json:"operStatus" param:"operStatus" query:"oper-status"`
	CreditCode   string `json:"creditCode" param:"creditCode" query:"credit-code"`
	Iden         string `json:"iden" param:"iden" query:"iden"`
	ApprDate     string `json:"apprDate" param:"apprDate" query:"appr-date"`
	CreateDate   string `json:"createDate" param:"createDate" query:"create-date"`
	Qual         string `json:"qual" param:"qual" query:"qual"`
	EntType      string `json:"entType" param:"entType" query:"ent-type"`
	RegAuthority string `json:"regAuthority" param:"regAuthority" query:"reg-authority"`
	EngName      string `json:"engName" param:"engName" query:"eng-name"`
	Addr         string `json:"addr" param:"addr" query:"addr"`
	SpecificAddr string `json:"specificAddr" param:"specificAddr" query:"specific-addr"`
	LegalPerson  string `json:"legalPerson" param:"legalPerson" query:"legal-person"`
}

type FullTextSearchQuery struct {
	Text string `json:"text" param:"text" query:"text"`
}

type CompanyQueryResult struct {
	Data       []*Company `json:"data"`
	TotalCount int64      `json:"totalCount"`
	PageNum    int        `json:"pageNum"`
	PageSize   int        `json:"pageSize"`
}

type SearchResult struct {
	Total int64
	Hits  []map[string]interface{}
}

type AiQueryRequest struct {
	Prompt string `json:"prompt"`
}
