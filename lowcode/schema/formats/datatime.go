package formats

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/timeutils"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"time"
)

// DateTimeFormat represents a datetime format.
var DateTimeFormat = &jsonschema.Format{
	Name: "date-time",
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

func ParseDateTime(v any) (time.Time, error) {
	return timeutils.AnyToTime(v, time.Now())
}

func ParseDate(v any) (time.Time, error) {
	date, err := timeutils.AnyToTime(v, time.Now())
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local), nil
}
