package fileutils

type ReadOptions struct {
	RootPath string
	WorkPath string
}

func newReadOption(o ...*ReadOptions) *ReadOptions {
	opts := &ReadOptions{}
	for _, o := range o {
		opts.RootPath = o.RootPath
		opts.WorkPath = o.WorkPath
	}
	return opts
}
