package context2

type Options struct {
	checkAuth *bool
	tenantId  *string
}

func NewOptions(opts ...Options) *Options {
	o := &Options{}
	for _, item := range opts {
		if item.checkAuth != nil {
			o.checkAuth = item.checkAuth
		}
		if item.tenantId != nil {
			o.tenantId = item.tenantId
		}
	}
	return o
}

func (o *Options) CheckAuth() bool {
	if o.checkAuth != nil {
		return *o.checkAuth
	}
	return true
}

func (o *Options) SetCheckAuth(val bool) *Options {
	o.checkAuth = &val
	return o
}
