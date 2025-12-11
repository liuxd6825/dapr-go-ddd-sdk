package query

type StatusFindByUserId struct {
	UserId string `json:"userId" param:"user-id" required:"true"`
}
