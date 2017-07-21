package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
    "html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kellegous/go/context"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	alpha  = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	prefix = ":"
)

var (
	errInvalidURL   = errors.New("Invalid URL")
	errRedirectLoop = errors.New(" I'm sorry, Dave. I'm afraid I can't do that")
)

// A very simple encoding of numeric ids. This is simply a base62 encoding
// prefixed with ":"
func encodeID(id uint64) string {
	n := uint64(len(alpha))
	b := make([]byte, 0, 8)
	if id == 0 {
		return "0"
	}

	b = append(b, ':')

	for id > 0 {
		b = append(b, alpha[id%n])
		id /= n
	}

	return string(b)
}

// Clean a shortcut name. Currently this just means stripping any leading
// ":" to avoid collisions with auto generated names.
func cleanName(name string) string {
	for strings.HasPrefix(name, prefix) {
		name = name[1:]
	}
	return name
}

// Parse the shortcut name from the give URL path, given the base URL that is
// handling the request.
func parseName(base, path string) string {
	t := path[len(base):]
	ix := strings.Index(t, "/")
	if ix == -1 {
		return t
	}
	return t[:ix]
}

// Used as an API response, this is a route with its associated shortcut name.
type routeWithName struct {
	Name string `json:"name"`
	*context.Route
}

// The response type for all API responses.
type msg struct {
	Ok    bool           `json:"ok"`
	Error string         `json:"error,omitempty"`
	Route *routeWithName `json:"route,omitempty"`
}

// Encode the given data to JSON and send it to the client.
func writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Panic(err)
	}
}

// Encode the given named route as a msg and send it to the client.
func writeJSONRoute(w http.ResponseWriter, name string, rt *context.Route) {
	writeJSON(w, &msg{
		Ok: true,
		Route: &routeWithName{
			Name:  name,
			Route: rt,
		},
	}, http.StatusOK)
}

// Encode a simple success msg and send it to the client.
func writeJSONOk(w http.ResponseWriter) {
	writeJSON(w, &msg{
		Ok: true,
	}, http.StatusOK)
}

// Encode an error response and send it to the client.
func writeJSONError(w http.ResponseWriter, err string) {
	writeJSON(w, &msg{
		Ok:    false,
		Error: err,
	}, http.StatusOK)
}

// Encode a generic backend error and send it to the client.
func writeJSONBackendError(w http.ResponseWriter, err error) {
	log.Printf("[error] %s", err)
	writeJSONError(w, "backend error")
}

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

// The handler that processes all API requests.
type apiHandler struct {
	ctx *context.Context
}

// Check that the given URL is suitable as a shortcut link.
func validateURL(r *http.Request, s string) error {
	u, err := url.Parse(s)
	if err != nil {
		return errInvalidURL
	}

	switch u.Scheme {
	case "http", "https", "mailto", "ftp":
		break
	default:
		return errInvalidURL
	}

	if r.Host == u.Host {
		return errRedirectLoop
	}

	return nil
}

// Handle a POST request to the API.
func apiPost(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	p := parseName("/api/url/", r.URL.Path)

	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "invalid json")
		return
	}

	// Handle delete requests
	if req.URL == "" {
		if p == "" {
			writeJSONError(w, "url required")
			return
		}

		if err := ctx.Del(p); err != nil {
			writeJSONBackendError(w, err)
			return
		}

		writeJSONOk(w)
		return
	}

	if err := validateURL(r, req.URL); err != nil {
		writeJSONError(w, err.Error())
		return
	}

	// If no name is specified, an ID must be generate.
	if p == "" {
		id, err := ctx.NextID()
		if err != nil {
			writeJSONBackendError(w, err)
			return
		}
		p = encodeID(id)
	}

	rt := context.Route{
		URL:  req.URL,
		Time: time.Now(),
	}

	if err := ctx.Put(p, &rt); err != nil {
		writeJSONBackendError(w, err)
		return
	}

	writeJSONRoute(w, p, &rt)
}

// Handle a GET request to the API.
func apiGet(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	p := parseName("/api/url/", r.URL.Path)

	if p == "" {
		writeJSONOk(w)
		return
	}

	rt, err := ctx.Get(p)
	if err == leveldb.ErrNotFound {
		writeJSONOk(w)
		return
	} else if err != nil {
		writeJSONBackendError(w, err)
		return
	}

	writeJSONRoute(w, p, rt)
}

func (h *apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		apiPost(h.ctx, w, r)
	case "GET":
		apiGet(h.ctx, w, r)
	default:
		writeJSONError(w, http.StatusText(http.StatusMethodNotAllowed))
	}
}

// The default handler responds to most requests. It is responsible for the
// shortcut redirects and for sending unmapped shortcuts to the edit page.
type defaultHandler struct {
	ctx *context.Context
}

func (h *defaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := parseName("/", r.URL.Path)
	if p == "" {
		http.Redirect(w, r, "/edit/", http.StatusTemporaryRedirect)
		return
	}

    rt, err := h.ctx.Get(p)
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
    t := template.New("links");
    contents, _ := linksHtmlBytes();
    t, err := t.Parse(string(contents));
    if err != nil {
        log.Printf("no template");
        return
    }

    routes, _ := h.ctx.GetAll();
    t.Execute(w, routes);
}

// Setup a Mux with all web routes.
func allRoutes(ctx *context.Context, admin bool, version string) *http.ServeMux {
    mux := http.NewServeMux()
    mux.Handle("/", &defaultHandler{ctx})
    mux.Handle("/api/url/", &apiHandler{ctx})
    mux.Handle("/links", &listLinksHandler{ctx});
    mux.HandleFunc("/edit/", func(w http.ResponseWriter, r *http.Request) {
        serveAsset(w, r, "index.html")
    })
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
    return mux
}

// ListenAndServe sets up all web routes, binds the port and handles incoming
// web requests.
func ListenAndServe(addr string, admin bool, version string, ctx *context.Context) error {
    return http.ListenAndServe(addr, allRoutes(ctx, admin, version))
}
