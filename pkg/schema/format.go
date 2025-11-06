package schema

import (
	times2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types/times"
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

func ParseDateTime(v any) (*times2.Date, error) {
	return times2.AsDate(v)
}

func ParseDate(v any) (*times2.Time, error) {
	date, err := times2.AsTime(v)
	return date, err
}
