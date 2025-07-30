package mapperutils

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"time"
)

// timeType
var timeType = reflect.TypeOf(time.Time{})

// MapToStructWithTag
func MapToStructWithTag(data map[string]interface{}, target interface{}, tagName string) error {
	dateLayouts := []string{
		time.RFC3339,          // z.B. "2024-07-30T08:55:00Z"
		"2006-01-02 15:04:05", // z.B. "2019-10-10 10:10:10"
		"2006-01-02",          // Nur Datum
	}
	timeHook := StringToTimeHookFuncMultiLayout(dateLayouts...)
	config := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			timeHook,
			mapstructure.StringToSliceHookFunc(","),
		),
		ErrorUnused: true,
		Result:      target,
		TagName:     tagName,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	if err := decoder.Decode(data); err != nil {
		return err
	}

	return nil
}

// StringToTimeHookFuncMultiLayout ist eine Hook-Funktions-Fabrik.
// Sie erstellt einen mapstructure-Hook, der versucht, einen String mit mehreren Layouts in ein time.Time-Objekt zu parsen.
func StringToTimeHookFuncMultiLayout(layouts ...string) mapstructure.DecodeHookFunc {
	// Die zurückgegebene Funktion ist der eigentliche Hook, den mapstructure verwenden wird.
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		// Der Hook ist nur für die Konvertierung von string zu time.Time relevant.
		// Wenn die Typen nicht übereinstimmen, geben Sie die Daten unverändert zurück,
		// damit mapstructure die Standardlogik anwenden kann.
		if f.Kind() != reflect.String || t != timeType {
			return data, nil
		}

		// Den String-Wert aus den Daten extrahieren.
		s, ok := data.(string)
		if !ok {
			// Dies sollte nicht passieren, wenn f.Kind() == reflect.String, ist aber eine sichere Überprüfung.
			return data, nil
		}

		// Versuchen, den String mit jedem der bereitgestellten Layouts zu parsen.
		for _, layout := range layouts {
			parsedTime, err := time.Parse(layout, s)
			if err == nil {
				// Sobald ein Layout erfolgreich ist, das Ergebnis zurückgeben.
				return parsedTime, nil
			}
		}

		// Wenn keine Layouts übereinstimmen, einen Fehler zurückgeben.
		return nil, fmt.Errorf("konnte '%s' mit keinem der bereitgestellten Zeitlayouts parsen", s)
	}
}
