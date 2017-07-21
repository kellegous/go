package web

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/kellegous/go/context"
	"github.com/syndtr/goleveldb/leveldb"
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

// The default handler responds to most requests. It is responsible for the
// shortcut redirects and for sending unmapped shortcuts to the edit page.
func getDefault(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	p := parseName("/", r.URL.Path)
	if p == "" {
		http.Redirect(w, r, "/edit/", http.StatusTemporaryRedirect)
		return
	}

	rt, err := ctx.Get(p)
	if err == leveldb.ErrNotFound {
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

type listLinksHandler struct {
	ctx *context.Context
}

func (h *listLinksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := template.New("links")
	contents, _ := linksHtmlBytes()
	t, err := t.Parse(string(contents))
	if err != nil {
		log.Printf("no template")
		return
	}

	routes, _ := h.ctx.GetAll()
	t.Execute(w, routes)
}

// ListenAndServe sets up all web routes, binds the port and handles incoming
// web requests.
func ListenAndServe(addr string, admin bool, version string, ctx *context.Context) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/url/", func(w http.ResponseWriter, r *http.Request) {
		apiURL(ctx, w, r)
	})
	mux.HandleFunc("/api/urls/", func(w http.ResponseWriter, r *http.Request) {
		apiURLs(ctx, w, r)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		getDefault(ctx, w, r)
	})
	mux.HandleFunc("/edit/", func(w http.ResponseWriter, r *http.Request) {
		serveAsset(w, r, "index.html")
	})
	mux.Handle("/links", &listLinksHandler{ctx})
	mux.HandleFunc("/s/", func(w http.ResponseWriter, r *http.Request) {
		serveAsset(w, r, r.URL.Path[len("/s/"):])
	})
	mux.HandleFunc("/:version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, version)
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	if admin {
		mux.Handle("/admin/", &adminHandler{ctx})
	}

	return http.ListenAndServe(addr, mux)
}
