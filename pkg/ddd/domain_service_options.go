package ddd

type DoCommandOption interface {
	GetIsValidOnly() *bool
}

type doCommandOption struct {
	IsValidOnly *bool
}

type IsValidOnlyOptions interface {
	GetIsValidOnly() *bool
}

func NewDoCommandOption(isValidOnly bool) DoCommandOption {
	return &doCommandOption{
		IsValidOnly: &isValidOnly,
	}
}

func NewDoCommandOptionMerges(opts ...DoCommandOption) DoCommandOption {
	res := &doCommandOption{}
	res.Merges(opts...)
	return res
}

func (o *doCommandOption) Merges(opts ...DoCommandOption) {
	if o == nil {
		return
	}
	for _, opt := range opts {
		if opt != nil && opt.GetIsValidOnly() != nil {
			o.IsValidOnly = opt.GetIsValidOnly()
		}
	}
}

func (o *doCommandOption) GetIsValidOnly() *bool {
	return o.IsValidOnly
}
