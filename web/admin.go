package web

import (
	"net/http"

	"github.com/kellegous/go/context"
)

type adminHandler struct {
	ctx *context.Context
}

func adminGet(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	p := parseName("/admin/", r.URL.Path)

	if p == "" {
		writeJSONOk(w)
		return
	}

	if p == "dumps" {
		if golinks, err := ctx.GetAll(); err != nil {
			writeJSONBackendError(w, err)
			return
		} else {
			writeJSON(w, golinks, http.StatusOK)
		}
	}

}

func (h *adminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		adminGet(h.ctx, w, r)
	default:
		writeJSONError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusOK) // fix
	}
}
