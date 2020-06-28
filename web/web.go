package web

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"

	"github.com/stgarf/go-links/backend"
	"github.com/stgarf/go-links/internal"
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
	log.Printf("GET %s from %s", r.URL.Path, r.RemoteAddr)
	p, s := parseName("/", r.URL.Path)
	if p == "" {
		log.Println("Redirecting bare request to /edit/")
		http.Redirect(w, r, "/edit/", http.StatusTemporaryRedirect)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	log.Printf("%s", p)
	if p == ".hidden_adminz" {
		log.Println("Redirecting to admin handler")
		adminGet(backend, w, r)
		// http.Redirect(w, r, "/edit/", http.StatusTemporaryRedirect)
		return
	}

	rt, err := backend.Get(ctx, p)
	if errors.Is(err, internal.ErrRouteNotFound) {
		log.Printf("Not found, redirecting for creation: %s", r.URL.Path)
		http.Redirect(w, r,
			fmt.Sprintf("/edit/%s", cleanName(p)),
			http.StatusTemporaryRedirect)
		return
	} else if err != nil {
		log.Panic(err)
	}

	http.Redirect(w, r,
		rt.URL+s,
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

func getLinks2(backend backend.Backend, w http.ResponseWriter, r *http.Request) {
	t, err := templateFromAssetFn(links2Html)
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

	// Return requested keyword in json
	mux.HandleFunc("/api/url/", func(w http.ResponseWriter, r *http.Request) {
		apiURL(backend, host, w, r)
	})
	// Return all keywords in json
	mux.HandleFunc("/api/urls/", func(w http.ResponseWriter, r *http.Request) {
		apiURLs(backend, host, w, r)
	})
	// Serve the index page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		getDefault(backend, w, r)
	})
	// Serve favicon
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		serveAsset(w, r, "favicon.ico")
	})
	// Serve edit page
	mux.HandleFunc("/edit/", func(w http.ResponseWriter, r *http.Request) {
		p, _ := parseName("/edit/", r.URL.Path)

		// if this is a banned name, just redirect to the local URI. That'll show em.
		if isBannedName(p) {
			http.Redirect(w, r, fmt.Sprintf("/%s", p), http.StatusTemporaryRedirect)
			return
		}

		serveAsset(w, r, "edit.html")
	})
	// Serve all links page
	mux.HandleFunc("/links/", func(w http.ResponseWriter, r *http.Request) {
		getLinks(backend, w, r)
	})
	mux.HandleFunc("/links2/", func(w http.ResponseWriter, r *http.Request) {
		getLinks2(backend, w, r)
	})
	// Server static assets
	mux.HandleFunc("/s/", func(w http.ResponseWriter, r *http.Request) {
		serveAsset(w, r, r.URL.Path[len("/s/"):])
	})
	// Serve version string... TODO(sgarf): remove?
	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, version)
	})
	// Serve instructions!
	htmlString := "<html><body bgcolor='#393939' text='#9e9e9e'>start by visiting /edit/&lt;yourKeyword&gt;. message <a style='color:#03a9f4; text-decoration:none;' href='slack://channel?team=T02874Q7H&id=CF26YT1HT'>@garf</a> on slack for help.</body></html>"
	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, htmlString)
	})
	// Serve healthcheck endpoint
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	// TODO(knorton): Remove the admin handler.
	if admin {
		mux.Handle("/.hidden_adminz/", &adminHandler{backend})
	}

	return http.ListenAndServe(addr, mux)
}
