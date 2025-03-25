package jsonschema_ext

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"github.com/liuxd6825/jsonschema/v6"
)

// DateTimeFormat represents a datetime format.
var DateTimeFormat = &jsonschema.Format{
	Name: "datetime",
	Validate: func(v any) error {
		if v, ok := v.(string); ok {
			_, err := ParseDateTime(v)
			return err
		}
		return nil
	},
}

// DateFormat represents a datetime format.
var DateFormat = &jsonschema.Format{
	Name: "date",
	Validate: func(v any) error {
		if v, ok := v.(string); ok {
			_, err := ParseDate(v)
			return err
		}
		return nil
	},
}

func ParseDateTime(v any) (*times.Date, error) {
	return times.AsDate(v)
}

func ParseDate(v any) (*times.Time, error) {
	date, err := times.AsTime(v)
	return date, err
}
