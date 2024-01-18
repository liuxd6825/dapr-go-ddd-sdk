package ddd

import (
	"context"
	"encoding/json"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
)

func SaveSnapshot(ctx context.Context, tenantId string, aggregateType string, aggregateId string, eventStoreKey string) (resErr error) {
	defer func() {
		resErr = errors.GetRecoverError(resErr, recover())
	}()

	aggregate, err := NewAggregate(aggregateType)
	if err != nil {
		return err
	}

	req := &dapr.LoadEventsRequest{
		TenantId:    tenantId,
		AggregateId: aggregateId,
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
		err = json.Unmarshal(bytes, aggregate)
		if err != nil {
			return err
		}
	}
	records := *resp.EventRecords
	if records != nil && len(records) > snapshotEventsMinCount {
		sequenceNumber := uint64(0)
		for _, record := range *resp.EventRecords {
			sequenceNumber = record.SequenceNumber
			if err = CallEventHandler(ctx, aggregate, &record); err != nil {
				return err
			}
		}

		snapshot := &dapr.SaveSnapshotRequest{
			TenantId:         tenantId,
			AggregateData:    aggregate,
			AggregateId:      aggregateId,
			AggregateType:    aggregateType,
			AggregateVersion: aggregate.GetAggregateVersion(),
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
