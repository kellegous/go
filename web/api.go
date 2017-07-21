package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/kellegous/go/context"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	alpha = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
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

func apiURLPost(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
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

func apiURLGet(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
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

func apiURLDelete(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	p := parseName("/api/url/", r.URL.Path)

	if p == "" {
		writeJSONError(w, "name required")
		return
	}

	if err := ctx.Del(p); err == leveldb.ErrNotFound {
		writeJSONError(w, fmt.Sprintf("%s not found", p))
		return
	} else if err != nil {
		writeJSONBackendError(w, err)
		return
	}

	writeJSONOk(w)
}

func apiURLsGet(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	// TODO(knorton): This will allow enumeration of the routes.
	writeJSONError(w, http.StatusText(http.StatusNotImplemented))
}

func apiURL(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		apiURLPost(ctx, w, r)
	case "GET":
		apiURLGet(ctx, w, r)
	case "DELETE":
		apiURLDelete(ctx, w, r)
	default:
		writeJSONError(w, http.StatusText(http.StatusMethodNotAllowed))
	}
}

func apiURLs(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		apiURLsGet(ctx, w, r)
	default:
		writeJSONError(w, http.StatusText(http.StatusMethodNotAllowed))
	}
}

// Setup ...
func Setup(m *http.ServeMux, ctx *context.Context) {
	m.HandleFunc("/api/url/", func(w http.ResponseWriter, r *http.Request) {
		apiURL(ctx, w, r)
	})

	m.HandleFunc("/api/urls/", func(w http.ResponseWriter, r *http.Request) {
		apiURLs(ctx, w, r)
	})
}
