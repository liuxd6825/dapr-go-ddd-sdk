package livyserver_test

import (
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/spark/livyserver"
)

func TestLivyServer_CreateBatch(t *testing.T) {
	server := newServer()
	resp, err := server.CreateBatch(livyserver.CreateBatchRequest{})
	if err != nil {
		t.Fatalf("CreateBatch failed: %v", err)
	}
	t.Log(resp)
}

func TestLivyServer_GetBatch(t *testing.T) {
	server := newServer()
	resp, err := server.GetBatchId(9)
	if err != nil {
		t.Fatalf("GetTaskId failed: %v", err)
	}
	t.Log(resp)
}

func TestLivyServer_GetSessions(t *testing.T) {
	server := newServer()
	resp, err := server.GetSessions(0, 10)
	if err != nil {
		t.Fatalf("GetSessions failed: %v", err)
	}
	t.Log(resp)
}

func newServer() *livyserver.LivyServer {
	return livyserver.NewLivyServer("192.168.120.200", 8998)
}
