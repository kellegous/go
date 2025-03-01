package web

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"

	"github.com/kellegous/go/internal"
	"github.com/kellegous/go/internal/backend"
)

// The default handler responds to most requests. It is responsible for the
// shortcut redirects and for sending unmapped shortcuts to the edit page.
func getDefault(
	backend backend.Backend,
	assets http.Handler,
	w http.ResponseWriter,
	r *http.Request,
) {
	p := parseName("/", r.URL.Path)
	if p == "" {
		r.URL.Path = "/s/"
		assets.ServeHTTP(w, r)
		// http.Redirect(w, r, "/edit/", http.StatusTemporaryRedirect)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	rt, err := backend.Get(ctx, p)
	if errors.Is(err, internal.ErrRouteNotFound) {
		http.Redirect(w, r,
			fmt.Sprintf("/edit/%s", cleanName(p)),
			http.StatusTemporaryRedirect)
		return
	} else if err != nil {
		log.Panic(err)
	}

	http.Redirect(w, r,
		rt.URL,
		http.StatusTemporaryRedirect)

}

// ListenAndServe sets up all web routes, binds the port and handles incoming
// web requests.
func ListenAndServe(
	backend backend.Backend,
	assets http.Handler,
) error {
	addr := viper.GetString("addr")
	admin := viper.GetBool("admin")
	version := viper.GetString("version")
	host := viper.GetString("host")

	mux := http.NewServeMux()

	mux.HandleFunc("/api/url/", func(w http.ResponseWriter, r *http.Request) {
		apiURL(backend, host, w, r)
	})

	mux.HandleFunc("/api/urls/", func(w http.ResponseWriter, r *http.Request) {
		apiURLs(backend, host, w, r)
	})

	mux.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, struct {
			Host string `json:"host"`
		}{host}, http.StatusOK)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		getDefault(backend, assets, w, r)
	})

	mux.HandleFunc("/edit/", func(w http.ResponseWriter, r *http.Request) {
		p := parseName("/edit/", r.URL.Path)

		// if this is a banned name, just redirect to the local URI. That'll show em.
		if isBannedName(p) {
			http.Redirect(w, r, fmt.Sprintf("/%s", p), http.StatusTemporaryRedirect)
			return
		}

		r.URL.Path = "/s/edit/"
		assets.ServeHTTP(w, r)
	})

	mux.HandleFunc("/links/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	})

	mux.Handle("/s/", assets)

	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, version)
	})

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "👍")
	})

	// TODO(knorton): Remove the admin handler.
	if admin {
		mux.Handle("/admin/", &adminHandler{backend})
	}

	return http.ListenAndServe(addr, mux)
}
