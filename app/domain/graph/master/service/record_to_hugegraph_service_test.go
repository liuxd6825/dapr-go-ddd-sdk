package service

import (
	"context"
	"testing"
)

func TestRecordToHugeGraphService_Import(t *testing.T) {
	s := NewRecord2HugeGraphService()
	err := s.Import(context.Background(), "s3a://import-data/record/record.parquet", "1001")
	if err != nil {
		t.Error(err)
	}
}
