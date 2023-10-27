package web

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/kellegous/golinks/pkg/internal"
	"github.com/kellegous/golinks/pkg/store"
)

type Server struct {
	store             store.Store
	assetProxyBaseURL *url.URL
	addr              string
	host              string
	admin             bool
}

func NewServer(
	store store.Store,
	opts ...Option,
) (*Server, error) {
	s := &Server{
		store: store,
		// TODO(knorton): apply defaults
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		getDefault(s, w, r)
	})

	ah, err := assetsHandler(s)
	if err != nil {
		return err
	}
	mux.Handle("/ui/", ah)

	mux.HandleFunc("/edit", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui/edit/", http.StatusTemporaryRedirect)
	})
	mux.HandleFunc("/edit/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui/edit/", http.StatusTemporaryRedirect)
	})

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "👍")
	})

	// TODO(knorton): Remove the admin handler.
	if s.admin {
		mux.Handle("/admin/", &adminHandler{s.store})
	}

	// TODO(knorton): Rebuild this.
	return http.ListenAndServe(s.addr, mux)
}

// The default handler responds to most requests. It is responsible for the
// shortcut redirects and for sending unmapped shortcuts to the edit page.
func getDefault(s *Server, w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "/" {
		http.Redirect(w, r, "/ui/", http.StatusTemporaryRedirect)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	p := parseName("/", path)

	pattern, err := regexp.Compile(p)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	rt, err := s.store.Get(ctx, pattern)
	if errors.Is(err, internal.ErrRouteNotFound) {
		http.Redirect(w, r,
			fmt.Sprintf("/ui/edit/%s", cleanName(p)),
			http.StatusTemporaryRedirect)
		return
	} else if err != nil {
		log.Panic(err)
	}

	http.Redirect(w, r,
		rt.URL,
		http.StatusTemporaryRedirect)

}
