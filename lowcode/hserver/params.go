package hserver

import (
	"github.com/kataras/iris/v12"
)

type Params struct {
	ictx iris.Context
}

func NewParams(ictx iris.Context) *Params {
	return &Params{ictx: ictx}
}
