package restapp

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd"

type RegisterEventType struct {
	EventType string
	Version   string
	NewFunc   ddd.NewEventFunc
}

func (r *RegisterEventType) GetEventType() string {
	return r.EventType
}

func (r *RegisterEventType) GetVersion() string {
	return r.Version
}

func (r *RegisterEventType) GetNewFunc() ddd.NewEventFunc {
	return r.NewFunc
}
