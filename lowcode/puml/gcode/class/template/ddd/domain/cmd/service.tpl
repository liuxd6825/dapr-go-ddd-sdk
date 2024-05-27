package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/singleutils"
	"github.com/liuxu6825/duxm-graph-service/pkg/cmd-service/domain/graph/command"
	"github.com/liuxu6825/duxm-graph-service/pkg/cmd-service/domain/graph/model"
	base_service "github.com/liuxu6825/duxm-graph-service/pkg/cmd-service/infrastructure/base/domain/service"
)

// GraphCommandDomainService
// @Description:  案件分析图（人员/公司/合同/产品/账户） 命令领域服务
type GraphCommandDomainService struct {
	base_service.BaseCommandDomainService
}

// GetGraphCommandDomainService
// @Description: 获取单例领域服务
// @return service.GraphQueryDomainService
func GetGraphCommandDomainService() *GraphCommandDomainService {
	return singleutils.CreateObj[*GraphCommandDomainService](func() *GraphCommandDomainService {
		return newGraphCommandDomainService()
	})
}

// NewGraphCommandDomainService
// @Description: 创建领域服务
// @return *GraphCommandDomainService
func newGraphCommandDomainService() *GraphCommandDomainService {
	return &GraphCommandDomainService{}
}

// GraphSetData
// @Description: 添加节点命令
// @receiver s
// @param ctx 上下文
// @param cmd 添加节点命令
// @return *model.GraphAggregate
// @return error
func (s *GraphCommandDomainService) GraphSetData(ctx context.Context, cmd *command.GraphSetDataCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// PaymentNodeCreate
// @Description: 添加节点命令
// @receiver s
// @param ctx 上下文
// @param cmd 添加节点命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) PaymentNodeCreate(ctx context.Context, cmd *command.PaymentNodeCreateCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// PaymentNodeDelete
// @Description: 删除节点命令
// @receiver s
// @param ctx 上下文
// @param cmd 删除节点命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) PaymentNodeDelete(ctx context.Context, cmd *command.PaymentNodeDeleteCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// PaymentNodeUpdate
// @Description: 更新节点命令
// @receiver s
// @param ctx 上下文
// @param cmd 更新节点命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) PaymentNodeUpdate(ctx context.Context, cmd *command.PaymentNodeUpdateCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// AccountNodeCreate
// @Description: 添加节点命令
// @receiver s
// @param ctx 上下文
// @param cmd 添加节点命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) AccountNodeCreate(ctx context.Context, cmd *command.AccountNodeCreateCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// AccountNodeDelete
// @Description: 删除节点命令
// @receiver s
// @param ctx 上下文
// @param cmd 删除节点命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) AccountNodeDelete(ctx context.Context, cmd *command.AccountNodeDeleteCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// AccountNodeUpdate
// @Description: 更新节点命令
// @receiver s
// @param ctx 上下文
// @param cmd 更新节点命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) AccountNodeUpdate(ctx context.Context, cmd *command.AccountNodeUpdateCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// CompanyNodeCreate
// @Description: 添加节点命令
// @receiver s
// @param ctx 上下文
// @param cmd 添加节点命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) CompanyNodeCreate(ctx context.Context, cmd *command.CompanyNodeCreateCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// CompanyNodeDelete
// @Description: 删除节点命令
// @receiver s
// @param ctx 上下文
// @param cmd 删除节点命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) CompanyNodeDelete(ctx context.Context, cmd *command.CompanyNodeDeleteCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// CompanyNodeUpdate
// @Description: 更新节点命令
// @receiver s
// @param ctx 上下文
// @param cmd 更新节点命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) CompanyNodeUpdate(ctx context.Context, cmd *command.CompanyNodeUpdateCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// ContractNodeCreate
// @Description: 添加节点命令
// @receiver s
// @param ctx 上下文
// @param cmd 添加节点命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) ContractNodeCreate(ctx context.Context, cmd *command.ContractNodeCreateCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// ContractNodeDelete
// @Description: 删除节点命令
// @receiver s
// @param ctx 上下文
// @param cmd 删除节点命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) ContractNodeDelete(ctx context.Context, cmd *command.ContractNodeDeleteCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// ContractNodeUpdate
// @Description: 更新节点命令
// @receiver s
// @param ctx 上下文
// @param cmd 更新节点命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) ContractNodeUpdate(ctx context.Context, cmd *command.ContractNodeUpdateCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// GraphCreate
// @Description: 新建分析图命令
// @receiver s
// @param ctx 上下文
// @param cmd 新建分析图命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) GraphCreate(ctx context.Context, cmd *command.GraphCreateCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// GraphDelete
// @Description: 删除分析图命令
// @receiver s
// @param ctx 上下文
// @param cmd 删除分析图命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) GraphDelete(ctx context.Context, cmd *command.GraphDeleteCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// GraphUpdate
// @Description: 更新分析图命令
// @receiver s
// @param ctx 上下文
// @param cmd 更新分析图命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) GraphUpdate(ctx context.Context, cmd *command.GraphUpdateCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// HumanNodeCreate
// @Description: 添加节点命令
// @receiver s
// @param ctx 上下文
// @param cmd 添加节点命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) HumanNodeCreate(ctx context.Context, cmd *command.HumanNodeCreateCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// HumanNodeDelete
// @Description: 删除节点命令
// @receiver s
// @param ctx 上下文
// @param cmd 删除节点命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) HumanNodeDelete(ctx context.Context, cmd *command.HumanNodeDeleteCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// HumanNodeUpdate
// @Description: 更新节点命令
// @receiver s
// @param ctx 上下文
// @param cmd 更新节点命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) HumanNodeUpdate(ctx context.Context, cmd *command.HumanNodeUpdateCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// ProductNodeCreate
// @Description: 添加节点命令
// @receiver s
// @param ctx 上下文
// @param cmd 添加节点命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) ProductNodeCreate(ctx context.Context, cmd *command.ProductNodeCreateCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// ProductNodeDelete
// @Description: 删除节点命令
// @receiver s
// @param ctx 上下文
// @param cmd 删除节点命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) ProductNodeDelete(ctx context.Context, cmd *command.ProductNodeDeleteCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// ProductNodeUpdate
// @Description: 更新节点命令
// @receiver s
// @param ctx 上下文
// @param cmd 更新节点命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) ProductNodeUpdate(ctx context.Context, cmd *command.ProductNodeUpdateCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// RecordNodeCreate
// @Description: 添加节点命令
// @receiver s
// @param ctx 上下文
// @param cmd 添加节点命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) RecordNodeCreate(ctx context.Context, cmd *command.RecordCreateCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// RecordNodeDelete
// @Description: 删除节点命令
// @receiver s
// @param ctx 上下文
// @param cmd 删除节点命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) RecordNodeDelete(ctx context.Context, cmd *command.RecordDeleteCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// RecordNodeUpdate
// @Description: 更新节点命令
// @receiver s
// @param ctx 上下文
// @param cmd 更新节点命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) RecordNodeUpdate(ctx context.Context, cmd *command.RecordNodeUpdateCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// RelationCreate
// @Description: 创建关系命令
// @receiver s
// @param ctx 上下文
// @param cmd 创建关系命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) RelationCreate(ctx context.Context, cmd *command.RelationCreateCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// RelationDelete
// @Description: 删除关系命令
// @receiver s
// @param ctx 上下文
// @param cmd 删除关系命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) RelationDelete(ctx context.Context, cmd *command.RelationDeleteCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// RelationUpdate
// @Description: 更新关系命令
// @receiver s
// @param ctx 上下文
// @param cmd 更新关系命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) RelationUpdate(ctx context.Context, cmd *command.RelationUpdateCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// UserCreate
// @Description: 新建用户
// @receiver s
// @param ctx 上下文
// @param cmd 新建用户
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) UserCreate(ctx context.Context, cmd *command.UserCreateCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// UserDelete
// @Description: 删除分析图命令
// @receiver s
// @param ctx 上下文
// @param cmd 删除分析图命令
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) UserDelete(ctx context.Context, cmd *command.UserDeleteCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// UserUpdate
// @Description: 更新用户
// @receiver s
// @param ctx 上下文
// @param cmd 更新用户
// @return *model.GraphCommandDomainService
// @return error
func (s *GraphCommandDomainService) UserUpdate(ctx context.Context, cmd *command.UserUpdateCommand, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	return s.doCommand(ctx, cmd, func() error {
		return cmd.Validate()
	}, opts...)
}

// doCommand
// @Description:
// @receiver s
// @param ctx
// @param cmd
// @return *model.GraphAggregate
// @return error
func (s *GraphCommandDomainService) doCommand(ctx context.Context, cmd ddd.Command, validateFunc func() error, opts ...ddd.DoCommandOption) (*model.GraphAggregate, error) {
	option := ddd.NewDoCommandOptionMerges(opts...)

	// 进行业务检查
	if validateFunc != nil {
		if err := validateFunc(); err != nil {
			return nil, err
		}
	} else if err := cmd.Validate(); err != nil {
		return nil, err
	}

	// 如果只是业务检查，则不执行领域命令，
	validOnly := option.GetIsValidOnly()
	if (validOnly == nil && cmd.GetIsValidOnly()) || (validOnly != nil && *validOnly == true) {
		return nil, nil
	}

	// 新建聚合根对象
	agg := s.NewAggregate()

	// 如果领域命令执行时出错，则返回错误
	if err := ddd.ApplyCommand(ctx, agg, cmd); err != nil {
		return nil, err
	}

	return agg, nil
}

// GetAggregateById
// @Description: 获取聚合对象
// @receiver s
// @param ctx 上下文
// @param tenantId 租户id
// @param id 主键id
// @return *graph_model.GraphCommandDomainService  聚合对象
// @return bool 是否找到聚合根对象
// @return error 错误对象
func (s *GraphCommandDomainService) GetAggregateById(ctx context.Context, tenantId string, id string) (*model.GraphAggregate, bool, error) {
	agg := s.NewAggregate()
	_, ok, err := ddd.LoadAggregate(ctx, tenantId, id, agg)
	return agg, ok, err
}

// NewAggregate
// @Description: 新建聚合对象
// @receiver s
// @return *graph_model.GraphCommandDomainService 聚合对象
func (s *GraphCommandDomainService) NewAggregate() *model.GraphAggregate {
	return model.NewGraphAggregate()
}
