package ddd

type GetTenant interface {
	GetTenantId() string
}

type GetCommandId interface {
	GetCommandId() string
}

type GetAggregateId interface {
	GetAggregateId() AggregateId
}

type GetUpdateMask interface {
	GetUpdateMask() []string
}

type GetIsValidOnly interface {
	GetIsValidOnly() bool
}

type IsAggregateCreateCommand interface {
	IsAggregateCreateCommand()
}

type Command interface {
	GetCommandId
	GetTenant
	GetAggregateId
	GetIsValidOnly
	Verify
}

type CreateCommand interface {
	Command
}

type DeleteCommand interface {
	Command
}

type UpdateCommand interface {
	Command
	GetUpdateMask
}
