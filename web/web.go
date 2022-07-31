package web

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"

	"github.com/ctSkennerton/shortlinks/backend"
	"github.com/ctSkennerton/shortlinks/internal"
)

// Serve a bundled asset over HTTP.
func serveAsset(w http.ResponseWriter, r *http.Request, name string) {
	n, err := AssetInfo(name)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	a, err := Asset(name)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.ServeContent(w, r, n.Name(), n.ModTime(), bytes.NewReader(a))
}

func templateFromAssetFn(fn func() (*asset, error)) (*template.Template, error) {
	a, err := fn()
	if err != nil {
		return nil, err
	}

	t := template.New(a.info.Name())
	return t.Parse(string(a.bytes))
}

// The default handler responds to most requests. It is responsible for the
// shortcut redirects and for sending unmapped shortcuts to the edit page.
func getDefault(backend backend.Backend, w http.ResponseWriter, r *http.Request) {
	p := parseName("/", r.URL.Path)
	if p == "" {
		http.Redirect(w, r, "/about/", http.StatusTemporaryRedirect)
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

func getLinks(backend backend.Backend, w http.ResponseWriter, r *http.Request) {
	t, err := templateFromAssetFn(linksHtml)
	if err != nil {
		log.Panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	rts, err := backend.GetAll(ctx)
	if err != nil {
		log.Panic(err)
	}

	if err := t.Execute(w, rts); err != nil {
		log.Panic(err)
	}
}

// ListenAndServe sets up all web routes, binds the port and handles incoming
// web requests.
func ListenAndServe(backend backend.Backend) error {
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
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		getDefault(backend, w, r)
	})
	mux.HandleFunc("/edit/", func(w http.ResponseWriter, r *http.Request) {
		p := parseName("/edit/", r.URL.Path)

		// if this is a banned name, just redirect to the local URI. That'll show em.
		if isBannedName(p) {
			http.Redirect(w, r, fmt.Sprintf("/%s", p), http.StatusTemporaryRedirect)
			return
		}

		serveAsset(w, r, "edit.html")
	})
	mux.HandleFunc("/about/", func(w http.ResponseWriter, r *http.Request) {
		serveAsset(w, r, "about.html")
	})
	mux.HandleFunc("/links/", func(w http.ResponseWriter, r *http.Request) {
		getLinks(backend, w, r)
	})
	mux.HandleFunc("/s/", func(w http.ResponseWriter, r *http.Request) {
		serveAsset(w, r, r.URL.Path[len("/s/"):])
	})
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
