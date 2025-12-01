package code_generator

import "text/template"

type Options struct {
	PackageName string
	Template    *template.Template
	Values      any
}

func NewOptions(options ...*Options) *Options {
	if len(options) == 0 {
		return &Options{}
	}
	o := &Options{}
	for _, item := range options {
		if item == nil {
			continue
		}
		if item.Template != nil {
			o.Template = item.Template
		}
		if item.PackageName != "" {
			o.PackageName = item.PackageName
		}
		if item.Values != nil {
			o.Values = item.Values
		}
	}
	return o
}
