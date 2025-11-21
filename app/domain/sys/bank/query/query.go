package query

type FindByIdQuery struct {
	Id string `json:"id" param:"id" required:"true"`
}
