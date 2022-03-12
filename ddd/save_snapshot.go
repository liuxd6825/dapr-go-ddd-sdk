package ddd

type SaveSnapshotRequest struct {
	TenantId          string                 `json:"tenantId"`
	AggregateId       string                 `json:"AggregateId"`
	AggregateType     string                 `json:"aggregateType"`
	AggregateData     interface{}            `json:"aggregateData"`
	AggregateRevision string                 `json:"aggregateRevision"`
	Metadata          map[string]interface{} `json:"metadata"`
	SequenceNumber    int64                  `json:"sequenceNumber"`
}

type SaveSnapshotResponse struct {
}
