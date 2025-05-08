package pkg

import "github.com/dop251/goja"

type PackageOptions struct {
	WorkPath string
}

type Package interface {
	NewProxy(vm *goja.Runtime, opts ...*PackageOptions) *Proxy
}

func NewPackageOptions() *PackageOptions {
	o := &PackageOptions{}
	return o
}

func (o *PackageOptions) Merge(opts ...PackageOptions) *PackageOptions {
	for _, opt := range opts {
		o.WorkPath = opt.WorkPath
	}
	return o
}

func (o *PackageOptions) SetWorkPath(path string) *PackageOptions {
	o.WorkPath = path
	return o
}
