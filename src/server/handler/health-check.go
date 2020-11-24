package handler

import (
	"fmt"
	"net/http"
)

// HealthCheckHandler handles health-check
func HealthCheckHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "R API: %s %s", "r-go-pact", "local")
		},
	)
}
