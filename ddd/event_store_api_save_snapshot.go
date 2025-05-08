package ddd

import (
	"context"
	"encoding/json"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
)

func SaveSnapshot(ctx context.Context, tenantId string, aggType string, aggId string, eventStoreKey string) (resErr error) {
	defer func() {
		resErr = errors.GetRecoverError(resErr, recover())
	}()

	agg, err := NewAggregateByType(aggType)
	if err != nil {
		return err
	}

	req := &dapr.LoadEventsRequest{
		TenantId:    tenantId,
		AggregateId: aggId,
	}
	resp, err := LoadEvents(ctx, req, eventStoreKey)
	if err != nil {
		return err
	}
	if resp.Snapshot == nil && (resp.EventRecords == nil || len(*resp.EventRecords) == 0) {
		return err
	}

	if resp.Snapshot != nil {
		bytes, err := json.Marshal(resp.Snapshot.AggregateData)
		if err != nil {
			return err
		}
		err = json.Unmarshal(bytes, agg)
		if err != nil {
			return err
		}
	}
	records := *resp.EventRecords
	if records != nil && len(records) > snapshotEventsMinCount {
		sequenceNumber := uint64(0)
		for _, record := range *resp.EventRecords {
			sequenceNumber = record.SequenceNumber
			if err = CallEventHandler(ctx, agg, &record); err != nil {
				return err
			}
		}

		snapshot := &dapr.SaveSnapshotRequest{
			TenantId:         tenantId,
			AggregateData:    agg,
			AggregateId:      aggId,
			AggregateType:    aggType,
			AggregateVersion: agg.GetAggregateVersion(),
			SequenceNumber:   sequenceNumber,
		}
		eventStorage, err := GetEventStore(eventStoreKey)
		if err != nil {
			return err
		}
		_, err = eventStorage.SaveSnapshot(ctx, snapshot)
		if err != nil {
			return err
		}
	}
	return err
}
