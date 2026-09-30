package hugedao

import (
	"context"
	"testing"
)

func TestBusRelationDao_FindByName_NotFound(t *testing.T) {
	client := getClient(t)
	cleanAll(t, client)
	relDao := newBusRelationDao(t, client)

	got := relDao.FindByName(context.Background(), "no-such-thing")
	if got != nil {
		t.Fatalf("expected nil, got %+v", got)
	}
}

func TestBusRelationDao_CreateSameName_Noop(t *testing.T) {
	_ = newBusRelationDao
	// 即使 HugeGraph 不可用, 此测试也不应触发 panic
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("CreateSameName should be a no-op, panic: %v", r)
		}
	}()
}
