package restapi

import (
	"context"
	"encoding/json"
	"log"

	"github.com/dapr/go-sdk/service/common"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/notify/event"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/notify/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type NotifyApi struct {
	rootPath      string
	notifyService *service.NotifyService
}

func NewNotifyApi(rootPath string, env *env.Env) *NotifyApi {
	return &NotifyApi{
		rootPath:      rootPath,
		notifyService: service.NewNotifyService(dapr.GetDaprClient(), env),
	}
}

func (s *NotifyApi) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath+"/notify", "notify.NotifyApi", s)
	controller.Post("/subscribe", "Subscribe")
	controller.Post("/publish", "Publish")
	controller.GetData("/ws", "RegisterClient")
	return controller
}

// RegisterClient
// 场景：前端建立长连接，等待任务完成通知
func (s *NotifyApi) RegisterClient(ctx iris.Context) {
	userID := ctx.URLParam("user_id")
	if userID == "" {
		ctx.StatusCode(400)
		return
	}
	// 升级连接，该 Goroutine 将被长期占用，直到断开
	s.notifyService.RegisterClient(ctx.ResponseWriter(), ctx.Request(), userID)
}

// Subscribe Dapr 订阅入口
// 场景：Temporal -> Dapr -> 随机的一个 Backend 实例 -> 此函数
func (s *NotifyApi) Subscribe(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	var evt event.TaskEvent
	if err := json.Unmarshal(e.Data.([]byte), &evt); err != nil {
		return false, err
	}

	log.Printf("Dapr received task for User: %s. Broadcasting to all nodes...", evt.UserID)

	responsePayload := map[string]interface{}{
		"type":    "TASK_UPDATE",
		"payload": evt,
	}
	s.notifyService.BroadcastToAllNodes(ctx, evt.UserID, responsePayload)
	return false, nil
}

// Publish
// @Description: 发布消息到前端
// @receiver s
// @param ctx
// @param e
func (s *NotifyApi) Publish(ctx context.Context, evt *event.TaskEvent) error {
	return s.notifyService.Publish(ctx, evt)
}
