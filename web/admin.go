package web

import (
	"net/http"

	"github.com/HALtheWise/o-links/context"
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
		if links, err := ctx.GetAll(); err != nil {
			writeJSONBackendError(w, err)
			return
		} else {
			writeJSON(w, links, http.StatusOK)
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
