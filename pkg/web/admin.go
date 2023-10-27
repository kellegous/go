package web

import (
	"net/http"

	"github.com/kellegous/golinks/pkg/store"
)

type adminHandler struct {
	store store.Store
}

func adminGet(s store.Store, w http.ResponseWriter, r *http.Request) {
	// TODO(knorton): Fix this.
	// p := parseName("/admin/", r.URL.Path)

	// if p == "" {
	// 	writeJSONOk(w)
	// 	return
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	// defer cancel()

	// if p == "dumps" {
	// 	if golinks, err := store.GetAll(ctx); err != nil {
	// 		writeJSONBackendError(w, err)
	// 		return
	// 	} else {
	// 		writeJSON(w, golinks, http.StatusOK)
	// 	}
	// }

}

func (h *adminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		adminGet(h.store, w, r)
	default:
		writeJSONError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusOK) // fix
	}
}
