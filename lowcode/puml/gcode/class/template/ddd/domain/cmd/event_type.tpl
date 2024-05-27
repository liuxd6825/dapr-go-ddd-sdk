package event

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
)

type GraphEventType uint32

const (
	AccountNodeDeleteEventType GraphEventType = iota
	AccountNodeUpdateEventType
	RelationCreateEventType
	ContractNodeCreateEventType
	HumanNodeUpdateEventType
	GraphDeleteEventType
	ProductNodeDeleteEventType
	UserUpdateEventType
	CompanyNodeDeleteEventType
	CompanyNodeUpdateEventType
	RelationUpdateEventType
	GraphCreateEventType
	AccountNodeCreateEventType
	HumanNodeCreateEventType
	UserDeleteEventType
	HumanNodeDeleteEventType
	ContractNodeDeleteEventType
	RelationDeleteEventType
	GraphUpdateEventType
	ProductNodeUpdateEventType
	ContractNodeUpdateEventType
	ProductNodeCreateEventType
	CompanyNodeCreateEventType
	UserCreateEventType

	RecordCreateEventType
	RecordUpdateEventType
	RecordDeleteEventType

	PaymentNodeCreateEventType
	PaymentNodeUpdateEventType
	PaymentNodeDeleteEventType

	GraphSetDataEventType
)

//
// String()
// @Description: 转换成字符串
//
func (p GraphEventType) String() string {
	switch p {
	// GraphNode

	case GraphDeleteEventType:
		return "duxm-graph.GraphDeleteEvent"
	case GraphCreateEventType:
		return "duxm-graph.GraphCreateEvent"
	case GraphUpdateEventType:
		return "duxm-graph.GraphUpdateEvent"
	case GraphSetDataEventType:
		return "duxm-graph.GraphSetDataEvent"

	// AccountNode

	case AccountNodeCreateEventType:
		return "duxm-graph.AccountNodeCreateEvent"
	case AccountNodeDeleteEventType:
		return "duxm-graph.AccountNodeDeleteEvent"
	case AccountNodeUpdateEventType:
		return "duxm-graph.AccountNodeUpdateEvent"

	// HumanNode

	case HumanNodeDeleteEventType:
		return "duxm-graph.HumanNodeDeleteEvent"
	case HumanNodeCreateEventType:
		return "duxm-graph.HumanNodeCreateEvent"
	case HumanNodeUpdateEventType:
		return "duxm-graph.HumanNodeUpdateEvent"

	// RelationNode

	case RelationUpdateEventType:
		return "duxm-graph.RelationUpdateEvent"
	case RelationCreateEventType:
		return "duxm-graph.RelationCreateEvent"
	case RelationDeleteEventType:
		return "duxm-graph.RelationDeleteEvent"

	// ContractNode

	case ContractNodeCreateEventType:
		return "duxm-graph.ContractNodeCreateEvent"
	case ContractNodeDeleteEventType:
		return "duxm-graph.ContractNodeDeleteEvent"
	case ContractNodeUpdateEventType:
		return "duxm-graph.ContractNodeUpdateEvent"

	// ProductNode

	case ProductNodeUpdateEventType:
		return "duxm-graph.ProductNodeUpdateEvent"
	case ProductNodeDeleteEventType:
		return "duxm-graph.ProductNodeDeleteEvent"
	case ProductNodeCreateEventType:
		return "duxm-graph.ProductNodeCreateEvent"

	// User

	case UserUpdateEventType:
		return "duxm-graph.UserUpdateEvent"
	case UserDeleteEventType:
		return "duxm-graph.UserDeleteEvent"
	case UserCreateEventType:
		return "duxm-graph.UserCreateEvent"

	// CompanyNode

	case CompanyNodeDeleteEventType:
		return "duxm-graph.CompanyNodeDeleteEvent"
	case CompanyNodeUpdateEventType:
		return "duxm-graph.CompanyNodeUpdateEvent"
	case CompanyNodeCreateEventType:
		return "duxm-graph.CompanyNodeCreateEvent"

	//RecordNode

	case RecordCreateEventType:
		return "duxm-graph.RecordCreateEvent"
	case RecordUpdateEventType:
		return "duxm-graph.RecordUpdateEvent"
	case RecordDeleteEventType:
		return "duxm-graph.RecordDeleteEvent"

	// PaymentNode

	case PaymentNodeCreateEventType:
		return "duxm-graph.PaymentNodeCreateEvent"
	case PaymentNodeUpdateEventType:
		return "duxm-graph.PaymentNodeUpdateEvent"
	case PaymentNodeDeleteEventType:
		return "duxm-graph.PaymentNodeDeleteEvent"

	default:
		return "UNKNOWN"
	}
}

//
// GetRegisterEventTypes
// @Description: 获取聚合根注册事件类型
// @return []restapp.RegisterEventType
//
func GetRegisterEventTypes() []restapp.RegisterEventType {
	return []restapp.RegisterEventType{
		{
			EventType: GraphSetDataEventType.String(),
			Version:   NewGraphSetDataEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &GraphSetDataEvent{} },
		},

		// Record

		{
			EventType: RecordCreateEventType.String(),
			Version:   NewRecordCreateEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &RecordCreateEvent{} },
		},
		{
			EventType: RecordDeleteEventType.String(),
			Version:   NewRecordDeleteEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &RecordDeleteEvent{} },
		},
		{
			EventType: RecordUpdateEventType.String(),
			Version:   NewRecordUpdateEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &RecordUpdateEvent{} },
		},

		// AccountNode

		{
			EventType: AccountNodeCreateEventType.String(),
			Version:   NewAccountNodeCreateEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &AccountNodeCreateEvent{} },
		},
		{
			EventType: AccountNodeDeleteEventType.String(),
			Version:   NewAccountNodeDeleteEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &AccountNodeDeleteEvent{} },
		},
		{
			EventType: AccountNodeUpdateEventType.String(),
			Version:   NewAccountNodeUpdateEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &AccountNodeUpdateEvent{} },
		},

		// CompanyNode

		{
			EventType: CompanyNodeCreateEventType.String(),
			Version:   NewCompanyNodeCreateEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &CompanyNodeCreateEvent{} },
		},
		{
			EventType: CompanyNodeDeleteEventType.String(),
			Version:   NewCompanyNodeDeleteEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &CompanyNodeDeleteEvent{} },
		},
		{
			EventType: CompanyNodeUpdateEventType.String(),
			Version:   NewCompanyNodeUpdateEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &CompanyNodeUpdateEvent{} },
		},

		// ContractNode

		{
			EventType: ContractNodeCreateEventType.String(),
			Version:   NewContractNodeCreateEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &ContractNodeCreateEvent{} },
		},
		{
			EventType: ContractNodeDeleteEventType.String(),
			Version:   NewContractNodeDeleteEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &ContractNodeDeleteEvent{} },
		},
		{
			EventType: ContractNodeUpdateEventType.String(),
			Version:   NewContractNodeUpdateEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &ContractNodeUpdateEvent{} },
		},

		// Graph

		{
			EventType: GraphCreateEventType.String(),
			Version:   NewGraphCreateEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &GraphCreateEvent{} },
		},
		{
			EventType: GraphDeleteEventType.String(),
			Version:   NewGraphDeleteEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &GraphDeleteEvent{} },
		},
		{
			EventType: GraphUpdateEventType.String(),
			Version:   NewGraphUpdateEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &GraphUpdateEvent{} },
		},

		// HumanNode

		{
			EventType: HumanNodeCreateEventType.String(),
			Version:   NewHumanNodeCreateEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &HumanNodeCreateEvent{} },
		},
		{
			EventType: HumanNodeDeleteEventType.String(),
			Version:   NewHumanNodeDeleteEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &HumanNodeDeleteEvent{} },
		},
		{
			EventType: HumanNodeUpdateEventType.String(),
			Version:   NewHumanNodeUpdateEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &HumanNodeUpdateEvent{} },
		},

		// ProductNode

		{
			EventType: ProductNodeCreateEventType.String(),
			Version:   NewProductNodeCreateEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &ProductNodeCreateEvent{} },
		},
		{
			EventType: ProductNodeDeleteEventType.String(),
			Version:   NewProductNodeDeleteEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &ProductNodeDeleteEvent{} },
		},
		{
			EventType: ProductNodeUpdateEventType.String(),
			Version:   NewProductNodeUpdateEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &ProductNodeUpdateEvent{} },
		},

		// Relation

		{
			EventType: RelationCreateEventType.String(),
			Version:   NewRelationCreateEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &RelationCreateEvent{} },
		},
		{
			EventType: RelationDeleteEventType.String(),
			Version:   NewRelationDeleteEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &RelationDeleteEvent{} },
		},
		{
			EventType: RelationUpdateEventType.String(),
			Version:   NewRelationUpdateEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &RelationUpdateEvent{} },
		},

		// User

		{
			EventType: UserCreateEventType.String(),
			Version:   NewUserCreateEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &UserCreateEvent{} },
		},
		{
			EventType: UserDeleteEventType.String(),
			Version:   NewUserDeleteEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &UserDeleteEvent{} },
		},
		{
			EventType: UserUpdateEventType.String(),
			Version:   NewUserUpdateEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &UserUpdateEvent{} },
		},

		// PaymentNode

		{
			EventType: PaymentNodeCreateEventType.String(),
			Version:   NewPaymentNodeCreateEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &PaymentNodeCreateEvent{} },
		},
		{
			EventType: PaymentNodeDeleteEventType.String(),
			Version:   NewPaymentNodeDeleteEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &PaymentNodeDeleteEvent{} },
		},
		{
			EventType: PaymentNodeUpdateEventType.String(),
			Version:   NewPaymentNodeUpdateEvent("").GetEventVersion(),
			NewFunc:   func() interface{} { return &PaymentNodeUpdateEvent{} },
		},
	}
}
