package v2

import (
	"net/http"

	"github.com/kellegous/golinks/pkg/backend"
)

func Setup(
	mux *http.ServeMux,
	be backend.Backend,
	host string,
) {
	mux.HandleFunc(
		"/api/v2/urls",
		func(w http.ResponseWriter, r *http.Request) {
		})

	// GET /api/v2/link/:id
	// POST /api/v2/link/:id
	// DELETE /api/v2/link/:id
	mux.HandleFunc(
		"/api/v2/link/",
		func(w http.ResponseWriter, r *http.Request) {
		})
}
