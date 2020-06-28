package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/stgarf/go-links/internal"
)

// Used as an API response, this is a route with its associated shortcut name.
type routeWithName struct {
	Name       string `json:"name"`
	SourceHost string `json:"source_host"`
	*internal.Route
}

// The response type for all API responses.
type msg struct {
	Ok bool `json:"ok"`
}

type msgErr struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

type msgRoute struct {
	Ok    bool           `json:"ok"`
	Route *routeWithName `json:"route"`
}

type msgRoutes struct {
	Ok     bool             `json:"ok"`
	Routes []*routeWithName `json:"routes"`
	Next   string           `json:"next"`
}

// Encode the given data to JSON and send it to the client.
func writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Panic(err)
	}
}

// Encode a simple success msg and send it to the client.
func writeJSONOk(w http.ResponseWriter) {
	writeJSON(w, &msg{
		Ok: true,
	}, http.StatusOK)
}

// Encode an error response and send it to the client.
func writeJSONError(w http.ResponseWriter, err string, status int) {
	writeJSON(w, &msgErr{
		Ok:    false,
		Error: err,
	}, status)
}

// Encode a generic backend error and send it to the client.
func writeJSONBackendError(w http.ResponseWriter, err error) {
	log.Printf("[error] %s", err)
	writeJSONError(w, "backend error", http.StatusInternalServerError)
}

// Encode the given named route as a msg and send it to the client.
func writeJSONRoute(w http.ResponseWriter, name string, rt *internal.Route, host string) {
	r := routeWithName{
		Name:  name,
		Route: rt,
	}

	if host != "" {
		r.SourceHost = host
	}

	writeJSON(w, &msgRoute{
		Ok:    true,
		Route: &r,
	}, http.StatusOK)
}
