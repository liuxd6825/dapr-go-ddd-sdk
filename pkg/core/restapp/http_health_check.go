package restapp

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dapr/go-sdk/service/common"
)

// AddHealthCheckHandler appends provided app health check handler.
func (s *HttpServer) AddHealthCheckHandler(route string, fn common.HealthCheckHandler) error {
	if fn == nil {
		return fmt.Errorf("health check handler required")
	}

	if !strings.HasPrefix(route, "/") {
		route = fmt.Sprintf("/%s", route)
	}

	s.app.HandleMany("ALL", route, optionsHandler(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if err := fn(r.Context()); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		})))
	return nil
}
