package server

type NotFoundError struct {
}

var (
	notFoundError = &NotFoundError{}
)

func NewExecutor(rctx *RContext) *Executor {
	return &Executor{rctx: rctx}
}

func (e *NotFoundError) Error() string {
	return "not found"
}
