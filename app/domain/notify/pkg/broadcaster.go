package pkg

import (
	"context"
	"encoding/json"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/redis/go-redis/v9"
)

// BroadcastMessage 内部广播消息结构
type BroadcastMessage struct {
	UserID  string      `json:"user_id"`
	Payload interface{} `json:"payload"`
}

var rdb *redis.Client

func InitRedis(env *env.Env) {
	rdb, ok := env.GetRedis("$notify")
	if ok && rdb != nil {
		// 启动监听协程
		go subscribeToBroadcast(rdb.Client)
	} else {
		panic("notify redis is nil")
	}
}

// BroadcastToAllNodes
// 1. 发送广播：收到 Dapr 消息的实例调用此方法
func BroadcastToAllNodes(userID string, payload interface{}) {
	msg := BroadcastMessage{UserID: userID, Payload: payload}
	bytes, _ := json.Marshal(msg)
	// 发布到 Redis "ws-broadcast" 频道
	rdb.Publish(context.Background(), "ws-broadcast", bytes)
}

// subscribeToBroadcast
// 2. 接收广播：所有实例都会运行此方法
func subscribeToBroadcast(rdb *redis.Client) {
	// 订阅 Redis 频道
	pubsub := rdb.Subscribe(context.Background(), "ws-broadcast")
	ch := pubsub.Channel()

	for msg := range ch {
		var data BroadcastMessage
		json.Unmarshal([]byte(msg.Payload), &data)

		// 关键点：尝试在本地推送。
		// 如果当前实例没有该用户的连接，LocalPush 会直接忽略，不会报错。
		// 如果当前实例有该用户的连接（1个或多个），都会收到通知。
		// Manager.LocalPush(data.UserID, data.Payload)
	}
}
