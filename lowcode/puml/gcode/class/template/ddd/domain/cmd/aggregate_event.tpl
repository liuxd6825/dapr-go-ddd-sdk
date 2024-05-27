package model

import (
	"context"
	"github.com/liuxu6825/duxm-graph-service/pkg/cmd-service/domain/graph/event"
	"github.com/liuxu6825/duxm-graph-service/pkg/cmd-service/infrastructure/utils"
)

//
// OnGraphSetDataEventV1s0
// @Description: GraphSetDataEventV1s0 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnGraphSetDataEventV1s0(ctx context.Context, e *event.GraphSetDataEvent) error {
	if len(e.Data.Id) == 0 {
		return nil
	}
	return nil
}

//
// OnAccountNodeCreateEventV1s0
// @Description: AccountNodeCreateEvent 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnAccountNodeCreateEventV1s0(ctx context.Context, e *event.AccountNodeCreateEvent) error {
	_, err := a.AccountNodes.AddMapper(ctx, e.Data.Id, &e.Data)
	return err
}

//
// OnAccountNodeDeleteEventV1s0
// @Description: AccountNodeDeleteEvent 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnAccountNodeDeleteEventV1s0(ctx context.Context, e *event.AccountNodeDeleteEvent) error {
	return a.AccountNodes.DeleteById(ctx, e.Data.Id)
}

//
// OnAccountNodeUpdateEventV1s0
// @Description: AccountNodeUpdateEvent 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnAccountNodeUpdateEventV1s0(ctx context.Context, e *event.AccountNodeUpdateEvent) error {
	_, _, err := a.AccountNodes.UpdateMapper(ctx, e.Data.Id, &e.Data, e.UpdateMask)
	return err
}

//
// OnCompanyNodeCreateEventV1s0
// @Description: CompanyNodeCreateEvent 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnCompanyNodeCreateEventV1s0(ctx context.Context, e *event.CompanyNodeCreateEvent) error {
	_, err := a.CompanyNodes.AddMapper(ctx, e.Data.Id, &e.Data)
	return err
}

//
// OnCompanyNodeDeleteEventV1s0
// @Description: CompanyNodeDeleteEvent 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnCompanyNodeDeleteEventV1s0(ctx context.Context, e *event.CompanyNodeDeleteEvent) error {
	return a.CompanyNodes.DeleteById(ctx, e.Data.Id)
}

//
// OnCompanyNodeUpdateEventV1s0
// @Description: CompanyNodeUpdateEvent 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnCompanyNodeUpdateEventV1s0(ctx context.Context, e *event.CompanyNodeUpdateEvent) error {
	_, _, err := a.CompanyNodes.UpdateMapper(ctx, e.Data.Id, &e.Data, e.UpdateMask)
	return err
}

//
// OnContractNodeCreateEventV1s0
// @Description: ContractNodeCreateEvent 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnContractNodeCreateEventV1s0(ctx context.Context, e *event.ContractNodeCreateEvent) error {
	_, err := a.ContractNodes.AddMapper(ctx, e.Data.Id, &e.Data)
	return err
}

//
// OnContractNodeDeleteEventV1s0
// @Description: ContractNodeDeleteEvent 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnContractNodeDeleteEventV1s0(ctx context.Context, e *event.ContractNodeDeleteEvent) error {
	return a.ContractNodes.DeleteById(ctx, e.Data.Id)
}

//
// OnContractNodeUpdateEventV1s0
// @Description: ContractNodeUpdateEvent 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnContractNodeUpdateEventV1s0(ctx context.Context, e *event.ContractNodeUpdateEvent) error {
	_, _, err := a.ContractNodes.UpdateMapper(ctx, e.Data.Id, &e.Data, e.UpdateMask)
	return err
}

//
// OnGraphCreateEventV1s0
// @Description: GraphCreateEvent 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnGraphCreateEventV1s0(ctx context.Context, e *event.GraphCreateEvent) error {
	return utils.Mapper(e.Data, a)
}

//
// OnGraphDeleteEventV1s0
// @Description: GraphDeleteEvent 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnGraphDeleteEventV1s0(ctx context.Context, e *event.GraphDeleteEvent) error {
	a.IsDeleted = true
	return nil
}

//
// OnGraphUpdateEventV1s0
// @Description: GraphUpdateEvent 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnGraphUpdateEventV1s0(ctx context.Context, e *event.GraphUpdateEvent) error {
	return utils.MaskMapperRemove(e.Data, a, e.UpdateMask, aggMapperRemove)
}

//
// OnHumanNodeCreateEventV1s0
// @Description: HumanNodeCreateEvent 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnHumanNodeCreateEventV1s0(ctx context.Context, e *event.HumanNodeCreateEvent) error {
	_, err := a.HumanNodes.AddMapper(ctx, e.Data.Id, &e.Data)
	return err
}

//
// OnHumanNodeDeleteEventV1s0
// @Description: HumanNodeDeleteEvent 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnHumanNodeDeleteEventV1s0(ctx context.Context, e *event.HumanNodeDeleteEvent) error {
	return a.HumanNodes.DeleteById(ctx, e.Data.Id)
}

//
// OnHumanNodeUpdateEventV1s0
// @Description: HumanNodeUpdateEvent 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnHumanNodeUpdateEventV1s0(ctx context.Context, e *event.HumanNodeUpdateEvent) error {
	_, _, err := a.HumanNodes.UpdateMapper(ctx, e.Data.Id, &e.Data, e.UpdateMask)
	return err
}

//
// OnProductNodeCreateEventV1s0
// @Description: ProductNodeCreateEvent 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnProductNodeCreateEventV1s0(ctx context.Context, e *event.ProductNodeCreateEvent) error {
	_, err := a.ProductNodes.AddMapper(ctx, e.Data.Id, &e.Data)
	return err
}

//
// OnProductNodeDeleteEventV1s0
// @Description: ProductNodeDeleteEvent 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnProductNodeDeleteEventV1s0(ctx context.Context, e *event.ProductNodeDeleteEvent) error {
	return a.ProductNodes.DeleteById(ctx, e.Data.Id)
}

//
// OnProductNodeUpdateEventV1s0
// @Description: ProductNodeUpdateEvent 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnProductNodeUpdateEventV1s0(ctx context.Context, e *event.ProductNodeUpdateEvent) error {
	_, _, err := a.ProductNodes.UpdateMapper(ctx, e.Data.Id, &e.Data, e.UpdateMask)
	return err
}

//
// OnRelationCreateEventV1s0
// @Description: RelationCreateEvent 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnRelationCreateEventV1s0(ctx context.Context, e *event.RelationCreateEvent) error {
	_, err := a.Relations.AddMapper(ctx, e.Data.Id, &e.Data)
	return err
}

//
// OnRelationDeleteEventV1s0
// @Description: RelationDeleteEvent 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnRelationDeleteEventV1s0(ctx context.Context, e *event.RelationDeleteEvent) error {
	return a.Relations.DeleteById(ctx, e.Data.Id)
}

//
// OnRelationUpdateEventV1s0
// @Description: RelationUpdateEvent 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnRelationUpdateEventV1s0(ctx context.Context, e *event.RelationUpdateEvent) error {
	_, _, err := a.Relations.UpdateMapper(ctx, e.Data.Id, &e.Data, e.UpdateMask)
	return err
}

//
// OnUserCreateEventV1s0
// @Description: UserCreateEvent 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnUserCreateEventV1s0(ctx context.Context, e *event.UserCreateEvent) error {
	_, err := a.Users.AddMapper(ctx, e.Data.Id, &e.Data)
	return err
}

//
// OnUserDeleteEventV1s0
// @Description: UserDeleteEvent 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnUserDeleteEventV1s0(ctx context.Context, e *event.UserDeleteEvent) error {
	return a.Users.DeleteById(ctx, e.Data.Id)
}

//
// OnUserUpdateEventV1s0
// @Description: UserUpdateEvent 领域事件 事件溯源处理器
// @receiver a
// @param ctx 上下文件
// @param event 领域事件
// @return err 错误
//
func (a *GraphAggregate) OnUserUpdateEventV1s0(ctx context.Context, e *event.UserUpdateEvent) error {
	_, _, err := a.Users.UpdateMapper(ctx, e.Data.Id, &e.Data, e.UpdateMask)
	return err
}

func (a *GraphAggregate) OnPaymentNodeCreateEventV1s0(ctx context.Context, e *event.PaymentNodeCreateEvent) {

}

func (a *GraphAggregate) OnPaymentNodeUpdateEventV1s0(ctx context.Context, e *event.PaymentNodeUpdateEvent) {

}

func (a *GraphAggregate) OnPaymentNodeDeleteEventV1s0(ctx context.Context, e *event.PaymentNodeDeleteEvent) {

}
