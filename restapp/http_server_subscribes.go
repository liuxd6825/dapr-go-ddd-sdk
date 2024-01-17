package restapp

import (
	"github.com/kataras/iris/v12/context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
)

// registerSubscribeHandler
// @Description: 新建领域事件控制器
// @param subscribes
// @param queryEventHandler
// @return ddd.SubscribeHandler
func (s *HttpServer) registerSubscribeHandler(subscribes []*ddd.Subscribe, queryEventHandler ddd.QueryEventHandler, interceptors []ddd.SubscribeInterceptorFunc) (ddd.SubscribeHandler, error) {

	subscribesHandler := func(sh ddd.SubscribeHandler, subscribe *ddd.Subscribe) (err error) {
		defer func() {
			err = errors.GetRecoverError(err, recover())
		}()
		s.app.Handle("POST", subscribe.Route, func(ictx *context.Context) {
			ctx, _ := NewContext(ictx)
			if err = sh.SubscribeHandler(ctx, ddd.NewSubscribeContext(ictx)); err != nil {
				ictx.SetErr(err)
			}
		})
		return err
	}

	handler := ddd.NewSubscribeHandler(subscribes, queryEventHandler, subscribesHandler, interceptors)
	if err := ddd.RegisterQueryHandler(handler, ddd.GetEventStoreDefaultPubsubName()); err != nil {
		return nil, err
	}
	return handler, nil
}

// subscribesHandler
//
//	@Description: Dapr获取订阅消息
//	@receiver s
//	@param ictx
func (s *HttpServer) subscribesHandler(ictx *context.Context) {
	ctx, _ := NewContext(ictx)
	defer func() {
		if err := errors.GetRecoverError(nil, recover()); err != nil {
			fields := logs.Fields{
				"func":  "restapp.HttpServer.subscribesHandler()",
				"error": err.Error(),
			}
			logs.Error(ctx, "", fields)
		}
	}()

	subscribes := ddd.GetSubscribes()
	ictx.Header("Context-Type", "application/json")
	_ = ictx.JSON(subscribes)

	if logs.GetLevel(ctx) >= logs.InfoLevel {
		for _, s := range subscribes {
			fields := logs.Fields{
				"dapr":   "subscribes",
				"pubsub": s.PubsubName,
				"topic":  s.Topic,
				"route":  s.Route,
			}
			logs.Info(ctx, "", fields)
		}
	}

}
