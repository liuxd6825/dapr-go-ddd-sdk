package dapr

import (
	"context"
	daprsdk "github.com/liuxd6825/dapr-go-sdk/client"
)

type Client = daprsdk.Client
type ClientOption = daprsdk.ClientOption

func newClient() (client Client, err error) {
	return daprsdk.NewClient()
}

// NewClientWithPort instantiates Dapr using specific gRPC port.
func newClientWithPort(port string) (client Client, err error) {
	return daprsdk.NewClientWithPort(port)
}

// NewClientWithAddress instantiates Dapr using specific address (including port).
// Deprecated: use NewClientWithAddressContext instead.
func newClientWithAddress(address string) (client Client, err error) {
	return daprsdk.NewClientWithAddress(address)
}

// NewClientWithAddressContext instantiates Dapr using specific address (including port).
// Uses the provided context to create the connection.
func newClientWithAddressContext(ctx context.Context, address string, opts ...ClientOption) (client Client, err error) {
	return daprsdk.NewClientWithAddressContext(ctx, address, opts...)
}

// SetMaxCallRecvMsgSize
//
//	@Description: 设置Client最大接收与发送数据包大小，单位M。需要在NewClient()之前设置。
//	@param val
func SetMaxCallRecvMsgSize(val int) {
	daprsdk.DefaultMaxCallRecvMsgSize = val
}

// GetMaxCallRecvMsgSize
//
//	@Description: 获取Client最大接收与发送数据包大小，单位M
//	@return int
func GetMaxCallRecvMsgSize() int {
	return daprsdk.DefaultMaxCallRecvMsgSize
}
