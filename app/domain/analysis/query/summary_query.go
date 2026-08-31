package query

type SummaryQuery struct {
	CaseId string `json:"case_id"`
	//MasterType string `json:"masterType"`
	//MasterId   string `json:"masterId"`
	DateMin        string   `json:"date_min"`
	DateMax        string   `json:"date_max"`
	Name           string   `json:"name"`
	OppName        string   `json:"opp_name"`
	TotalAmountMin *float64 `json:"amount_min"`
	TotalAmountMax *float64 `json:"amount_max"`
	TotalCountMin  *int     `json:"count_min"`
	TotalCountMax  *int     `json:"count_max"`
	Page           int      `json:"page"`
	PageSize       int      `json:"page_size"`
	Filter         string   `json:"filter"`
	Sort           string   `json:"sort"`
}
