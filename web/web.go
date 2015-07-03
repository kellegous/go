package web

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func cleanName(name string) string {
	for strings.HasPrefix(name, prefix) {
		name = name[1:]
	}
	return name
}

func parseName(base, path string) string {
	t := path[len(base):]
	ix := strings.Index(t, "/")
	if ix == -1 {
		return t
	}
	return t[:ix]
}

func writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Panic(err)
	}
}

func writeJSONError(w http.ResponseWriter, error string, status int) {
	writeJSON(w, map[string]interface{}{
		"error": error,
	}, status)
}

func writeJSONNotFound(w http.ResponseWriter) {
	writeJSON(w, nil, http.StatusNotFound)
}

func writeJSONRoute(w http.ResponseWriter, name string, rt *context.Route) {
	res := struct {
		Name string    `json:"name"`
		URL  string    `json:"url"`
		Time time.Time `json:"time"`
	}{
		name,
		rt.URL,
		rt.Time,
	}

	writeJSON(w, &res, http.StatusOK)
}

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

type apiHandler struct {
	ctx *context.Context
}

func validURL(s string) bool {
	u, err := url.Parse(s)
	if err != nil {
		return false
	}

	switch u.Scheme {
	case "http", "https", "mailto", "ftp":
		return true
	}

	return false
}

func apiPost(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	p := parseName("/api/url/", r.URL.Path)

	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Handle delete requests
	if req.URL == "" {
		if p == "" {
			writeJSONError(w, "url required", http.StatusBadRequest)
			return
		}

		if err := ctx.Del(p); err != nil {
			log.Panic(err)
		}
		return
	}

	if !validURL(req.URL) {
		writeJSONError(w, "invalid URL", http.StatusBadRequest)
		return
	}

	// If no name is specified, an ID must be generate.
	if p == "" {
		id, err := ctx.NextID()
		if err != nil {
			log.Panic(err)
		}
		p = encodeID(id)
	}

	rt := context.Route{
		URL:  req.URL,
		Time: time.Now(),
	}

	if err := ctx.Put(p, &rt); err != nil {
		log.Panic(err)
	}

	writeJSONRoute(w, p, &rt)
}

func apiGet(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	p := parseName("/api/url/", r.URL.Path)

	if p == "" {
		writeJSONNotFound(w)
		return
	}

	rt, err := ctx.Get(p)
	if err == leveldb.ErrNotFound {
		writeJSONNotFound(w)
		return
	} else if err != nil {
		log.Panic(err)
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
		writeJSONError(w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed)
	}
}

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

type editHandler struct {
	ctx *context.Context
}

func (h *editHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serveAsset(w, r, "index.html")
}

// ListenAndServe ...
func ListenAndServe(addr string, ctx *context.Context) error {
	mux := http.NewServeMux()

	mux.Handle("/", &defaultHandler{ctx})
	mux.Handle("/edit/", &editHandler{ctx})
	mux.Handle("/api/url/", &apiHandler{ctx})
	mux.HandleFunc("/s/", func(w http.ResponseWriter, r *http.Request) {
		serveAsset(w, r, r.URL.Path[len("/s/"):])
	})

	return http.ListenAndServe(addr, mux)
}
