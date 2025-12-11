package restapi

import (
	"context"

	"log"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/notify/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/notify/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type NotifyApi struct {
	rootPath      string
	notifyService *service.NotifyService
	websocket     *neffos.Server
}

func NewNotifyApi(rootPath string, env *env.Env) *NotifyApi {
	return &NotifyApi{
		rootPath:      rootPath,
		notifyService: service.NewNotifyService(dapr.GetDaprClient(), env),
	}
}

func (s *NotifyApi) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath+"/sys/notify", "sys.notify", s)
	//controller.EventHandle("/subs-message-event", "SubsNotify")
	//controller.EventHandle("/subs-status-event", "SubsActivity")
	controller.Post("/message", "PublishMessage")
	controller.Post("/status", "PublishStatus")
	s.initWebsocket(app)
	return controller
}

func (s *NotifyApi) initWebsocket(app *iris.Application) {
	// 1. 配置 Websocket 服务器
	ws := websocket.New(websocket.DefaultGorillaUpgrader, neffos.Events{})

	ws.OnConnect = func(c *neffos.Conn) error {
		log.Printf("connect to websocket")
		return nil
	}

	ws.OnConnect = func(c *websocket.Conn) error {
		log.Printf("[%s] Connected to server! ", c.ID())
		return nil
	}

	ws.OnDisconnect = func(c *websocket.Conn) {
		log.Printf("[%s] Disconnected from server ", c.ID())
	}

	ws.OnUpgradeError = func(err error) {
		log.Printf("Upgrade Error: %s", err.Error())
	}

	// 2. 注册 WebSocket 路由
	// Endpoint: /ws
	app.Get("/api/v1.0/notify/ws", websocket.Handler(ws))

}

// Register
// 场景：前端建立长连接，等待任务完成通知
func (s *NotifyApi) Register(ictx iris.Context) {
	userID := ictx.URLParam("user_id")
	if userID == "" {
		ictx.StatusCode(400)
		return
	}
	// 升级连接，该 Goroutine 将被长期占用，直到断开
	s.notifyService.RegisterClient(ictx.ResponseWriter(), ictx.Request(), userID)
}

/*func (s *NotifyApi) SubsNotify(ctx context.Context, e *model.Message) (retry bool, err error) {
	log.Printf("Dapr received task for User: %s. Broadcasting to all nodes...", e.Id)
	s.notifyService.PublishMessage(ctx, e)
	return false, nil
}

func (s *NotifyApi) SubsStatus(ctx context.Context, e *model.Status) (retry bool, err error) {
	log.Printf("Dapr received task for User: %s. Broadcasting to all nodes...", e.Id)
	s.notifyService.PublishMessage(ctx, e)
	return false, nil
}
*/

// PublishMessage
// @Description: 发布消息到前端
// @receiver s
// @param ctx
// @param e
func (s *NotifyApi) PublishMessage(ctx context.Context, evt *model.Message) {
	s.notifyService.PublishMessage(ctx, evt)
}

func (s *NotifyApi) PublishStatus(ctx context.Context, evt *model.Status) {
	s.notifyService.PublishStatus(ctx, evt)
}
