package runtime

type Options struct {
	WorkPath string
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) SetWorkPath(path string) *Options {
	o.WorkPath = path
	return o
}

func (o *Options) GetWorkPath() string {
	return o.WorkPath
}
