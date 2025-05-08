package ddd

type GetTenantId interface {
	GetTenantId() string
}
type SetTenantId interface {
	SetTenantId(string)
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
	GetTenantId
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
