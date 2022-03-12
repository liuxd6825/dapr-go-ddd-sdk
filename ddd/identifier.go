package ddd

type Identifier struct {
	AggregateId string `json:"aggregateId"`
	DetailId    string `json:"detailId"`
}

func NewIdentifier(aggregateId string, detailId string) *Identifier {
	return &Identifier{
		AggregateId: aggregateId,
		DetailId:    detailId,
	}
}

func (i *Identifier) GetAggregateId() string {
	return i.AggregateId
}

func (i *Identifier) GetDetailId() string {
	return i.DetailId
}

func (i *Identifier) IsDetail() bool {
	if len(i.DetailId) == 0 {
		return false
	}
	return true
}
