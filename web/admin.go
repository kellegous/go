package web

import (
	"context"
	"net/http"
	"time"

	"github.com/stgarf/go-links/backend"
)

type adminHandler struct {
	backend backend.Backend
}

func adminGet(backend backend.Backend, w http.ResponseWriter, r *http.Request) {
	p, _ := parseName("/.hidden_adminz/", r.URL.Path)

	if p == "" {
		writeJSONOk(w)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if p == "dumps" {
		if golinks, err := backend.GetAll(ctx); err != nil {
			writeJSONBackendError(w, err)
		} else {
			writeJSON(w, golinks, http.StatusOK)
		}
	}
}

func (h *adminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		adminGet(h.backend, w, r)
	default:
		writeJSONError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusOK) // fix
	}
}
