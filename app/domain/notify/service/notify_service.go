package service

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/notify/event"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/notify/pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/redis/go-redis/v9"
)

type BroadcastMessage struct {
	UserID  string      `json:"user_id"`
	Payload interface{} `json:"payload"`
}

type NotifyService struct {
	daprClient  dapr.Client
	redisClient *redis.Client
	env         *env.Env
	websockets  *pkg.WebsocketManager
}

var notifyServiceOnce sync.Once
var notifyServiceInstance *NotifyService

func NewNotifyService(client dapr.Client, env *env.Env) *NotifyService {
	notifyServiceOnce.Do(func() {
		notifyServiceInstance = newNotifyService(client, env)
	})
	return notifyServiceInstance
}

func newNotifyService(client dapr.Client, env *env.Env) *NotifyService {
	redisClient, ok := env.GetRedis(config.GetRedisKey())
	if !ok {
		panic("redis client not exist")
	}
	service := &NotifyService{
		daprClient:  client,
		redisClient: redisClient.Client,
		env:         env,
		websockets:  pkg.NewWebsocketManager(),
	}
	return service
}

// Publish
// 在 Temporal Activity 中调用此方法
func (n *NotifyService) Publish(ctx context.Context, evt *event.TaskEvent) error {
	// 使用 Dapr 发布消息
	return n.daprClient.PublishEvent(ctx, event.PubSubName, event.TopicName, evt)
}

// Close
// @Description: 关闭服务
// @receiver n
func (n *NotifyService) Close() {
	n.daprClient.Close()
	_ = n.redisClient.Close()
}

// BroadcastToAllNodes
// @Description:  发送广播：收到 Dapr 消息的实例调用此方法
func (n *NotifyService) BroadcastToAllNodes(ctx context.Context, userID string, payload interface{}) {
	msg := BroadcastMessage{UserID: userID, Payload: payload}
	bytes, _ := json.Marshal(msg)
	// 发布到 Redis "ws-broadcast" 频道
	n.redisClient.Publish(context.Background(), "ws-broadcast", bytes)
}

func (n *NotifyService) RegisterClient(w http.ResponseWriter, r *http.Request, userID string) {
	n.websockets.RegisterClient(w, r, userID)
}

// subscribeToBroadcast
// 2. 接收广播：所有实例都会运行此方法
func (n *NotifyService) subscribeToBroadcast(rdb *redis.Client) {
	// 订阅 Redis 频道
	pubsub := rdb.Subscribe(context.Background(), "ws-broadcast")
	ch := pubsub.Channel()

	for msg := range ch {
		var data BroadcastMessage
		json.Unmarshal([]byte(msg.Payload), &data)

		// 关键点：尝试在本地推送。
		// 如果当前实例没有该用户的连接，LocalPush 会直接忽略，不会报错。
		// 如果当前实例有该用户的连接（1个或多个），都会收到通知。
		n.websockets.LocalPush(data.UserID, data.Payload)
	}
}
