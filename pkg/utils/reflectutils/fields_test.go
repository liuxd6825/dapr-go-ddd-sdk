package reflectutils

import (
	"time"
)

type Super struct {
	F1 string
	F2 int
	F3 time.Time
}

type Object struct {
	Super
	Str string
	Int int64
}
