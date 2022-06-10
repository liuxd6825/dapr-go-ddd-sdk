package ddd_domain_service

type DoOptions struct {
	IsValidOnly *bool
}

type IsValidOnlyOptions interface {
	GetIsValidOnly() *bool
}

func NewDoOptions(isValidOnly bool) *DoOptions {
	return &DoOptions{
		IsValidOnly: &isValidOnly,
	}
}

func NewDoOptionsEmpty() *DoOptions {
	return &DoOptions{}
}

func NewDoOptionsMerges(opts ...*DoOptions) *DoOptions {
	res := NewDoOptions(false)
	res.Merges(opts...)
	return res
}

func (o *DoOptions) Merges(opts ...*DoOptions) {
	if o == nil {
		return
	}
	for _, opt := range opts {
		if opt != nil && opt.IsValidOnly != nil {
			o.IsValidOnly = opt.IsValidOnly
		}
	}
}

func (o *DoOptions) GetIsValidOnly() *bool {
	return o.IsValidOnly
}
