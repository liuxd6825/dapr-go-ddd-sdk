package eventbus_pkg

import "time"

type Aggregate struct {
	AggId       string    `json:"aggId"`
	AggType     string    `json:"aggType"`
	AggVer      string    `json:"aggVer"`
	CreatedTime time.Time `json:"createdTime"`
	TenantId    string    `json:"tenantId"`
}

func (a *Aggregate) GetAggId() string {
	return a.AggId
}

func (a *Aggregate) GetAggType() string {
	return a.AggType
}

func (a *Aggregate) GetAggVer() string {
	return a.AggVer
}

func (a *Aggregate) GetTenantId() string {
	return a.TenantId
}
