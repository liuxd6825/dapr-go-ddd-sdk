package light_rag

import (
	"context"
	"testing"
)

func Test_Query(t *testing.T) {
	ctx := context.Background()
	query := NewQueryRequest("孙悟空有手机吗", ModelLocal)
	resp, err := client().Query().Query(ctx, query)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(resp.Response)
}

func Test_QueryStream(t *testing.T) {
	ctx := context.Background()
	query := NewQueryRequest("孙悟空有手机是怎么来的", ModelLocal)
	err := client().Query().QueryStream(ctx, query, func(resp *QueryResponse) {
		t.Log(resp.Response)
	})
	if err != nil {
		t.Fatal(err)
	}
}
