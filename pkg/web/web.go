package web

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"

	"github.com/kellegous/golinks/pkg/backend"
	"github.com/kellegous/golinks/pkg/internal"
)

// The default handler responds to most requests. It is responsible for the
// shortcut redirects and for sending unmapped shortcuts to the edit page.
func getDefault(backend backend.Backend, w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "/" {
		http.Redirect(w, r, "/ui/", http.StatusTemporaryRedirect)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	p := parseName("/", path)

	rt, err := backend.Get(ctx, p)
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

func getLinks(backend backend.Backend, w http.ResponseWriter, r *http.Request) {
	// TODO(knorton): Restore this.

	// t, err := templateFromAssetFn(linksHtml)
	// if err != nil {
	// 	log.Panic(err)
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	// defer cancel()

	// rts, err := backend.GetAll(ctx)
	// if err != nil {
	// 	log.Panic(err)
	// }

	// if err := t.Execute(w, rts); err != nil {
	// 	log.Panic(err)
	// }
}

// ListenAndServe sets up all web routes, binds the port and handles incoming
// web requests.
func ListenAndServe(be backend.Backend, opts ...Option) error {
	var options Options
	for _, opt := range opts {
		if err := opt(&options); err != nil {
			return err
		}
	}

	addr := viper.GetString("addr")
	admin := viper.GetBool("admin")
	version := viper.GetString("version")
	host := viper.GetString("host")

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		getDefault(be, w, r)
	})

	ah, err := assetsHandler(&options)
	if err != nil {
		return err
	}
	mux.Handle("/ui/", ah)

	mux.HandleFunc("/api/url/", func(w http.ResponseWriter, r *http.Request) {
		apiURL(be, host, w, r)
	})
	mux.HandleFunc("/api/urls/", func(w http.ResponseWriter, r *http.Request) {
		apiURLs(be, host, w, r)
	})

	// added to skip a redirect from /edit -> /edit/
	mux.HandleFunc("/edit", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui/edit/", http.StatusTemporaryRedirect)
	})
	mux.HandleFunc("/edit/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui/edit/", http.StatusTemporaryRedirect)
	})

	mux.HandleFunc("/links/", func(w http.ResponseWriter, r *http.Request) {
		getLinks(be, w, r)
	})
	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, version)
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "👍")
	})

	// TODO(knorton): Remove the admin handler.
	if admin {
		mux.Handle("/admin/", &adminHandler{be})
	}

	return http.ListenAndServe(addr, mux)
}
