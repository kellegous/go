package web

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/kellegous/glue/logging"
	"go.uber.org/zap"

	"github.com/kellegous/golinks/internal"
	"github.com/kellegous/golinks/internal/store"
)

func getDefaultHandler(
	s store.Store,
	assets http.Handler,
) http.Handler {
	assets = http.StripPrefix("/s", assets)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		prefix := getPrefix(r)
		if prefix == "" {
			// this is the / route, serve /s/index.html
			assets.ServeHTTP(w, r)
			return
		}

		lg := logging.L(r.Context())

		proto, err := s.Get(r.Context(), prefix)
		if errors.Is(err, store.ErrLinkNotfound) {
			http.Redirect(w, r, fmt.Sprintf("/edit/%s", prefix), http.StatusTemporaryRedirect)
			return
		} else if err != nil {
			lg.Panic("failed to get link", zap.String("prefix", prefix), zap.Error(err))
			return
		}

		link, err := internal.ToLink(proto)
		if err != nil {
			lg.Panic("failed to convert proto to link", zap.String("prefix", prefix), zap.Error(err))
			return
		}

		if eu := link.Expand(r.URL.Path); eu != nil {
			http.Redirect(w, r, eu.URL, http.StatusTemporaryRedirect)
			return
		}

		http.NotFound(w, r)
	})
}

func getPrefix(r *http.Request) string {
	path := strings.TrimPrefix(r.URL.Path, "/")
	p, _, _ := strings.Cut(path, "/")
	return p
}
