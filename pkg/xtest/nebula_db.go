package xtest

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"

const (
	NebulaDBKey      = "default"
	NebulaHostLocal  = "localhost"
	NebulaHostRemote = "192.168.120.224"
)

func GetNebulaRemoteConfig() *env.Nebula {
	return &env.Nebula{
		Name:     NebulaDBKey,
		Addrs:    []string{NebulaHostRemote + ":9669"},
		User:     "root",
		Password: "nebula",
		Space:    "test_default",
		PoolSize: 10,
	}
}