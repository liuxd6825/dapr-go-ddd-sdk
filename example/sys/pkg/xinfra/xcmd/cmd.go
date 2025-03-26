package xcmd

type Cmd[T any] struct {
	CommandId string `json:"commandId"`
	EventType string `json:"eventType"`
	Data      T      `json:"data"`
}
