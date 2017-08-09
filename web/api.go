package web

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kellegous/go/context"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	alpha = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	errInvalidURL        = errors.New("Invalid URL")
	errRedirectLoop      = errors.New(" I'm sorry, Dave. I'm afraid I can't do that")
	genURLPrefix    byte = ':'
	postGenCursor        = []byte{genURLPrefix + 1}
)

// A very simple encoding of numeric ids. This is simply a base62 encoding
// prefixed with ":"
func encodeID(id uint64) string {
	n := uint64(len(alpha))
	b := make([]byte, 0, 8)
	if id == 0 {
		return "0"
	}

	b = append(b, genURLPrefix)

	for id > 0 {
		b = append(b, alpha[id%n])
		id /= n
	}

	return string(b)
}

// Advance to the next contetxt id and encode it as an ID.
func nextEncodedID(ctx *context.Context) (string, error) {
	id, err := ctx.NextID()
	if err != nil {
		return "", err
	}
	return encodeID(id), nil
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
		writeJSONError(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		writeJSONError(w, "url required", http.StatusBadRequest)
		return
	}

	if isBannedName(p) {
		writeJSONError(w, "name cannot be used", http.StatusBadRequest)
		return
	}

	if err := validateURL(r, req.URL); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If no name is specified, an ID must be generate.
	if p == "" {
		var err error
		p, err = nextEncodedID(ctx)
		if err != nil {
			writeJSONBackendError(w, err)
			return
		}
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
		writeJSONError(w, "no name given", http.StatusBadRequest)
		return
	}

	rt, err := ctx.Get(p)
	if err == leveldb.ErrNotFound {
		writeJSONError(w, "Not Found", http.StatusNotFound)
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
		writeJSONError(w, "name required", http.StatusBadRequest)
		return
	}

	if err := ctx.Del(p); err != nil {
		writeJSONBackendError(w, err)
		return
	}

	writeJSONOk(w)
}

func parseCursor(v string) ([]byte, error) {
	if v == "" {
		return nil, nil
	}

	return base64.URLEncoding.DecodeString(v)
}

func parseInt(v string, def int) (int, error) {
	if v == "" {
		return def, nil
	}

	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, err
	}

	return int(i), nil
}

func parseBool(v string, def bool) (bool, error) {
	if v == "" {
		return def, nil
	}

	v = strings.ToLower(v)
	if v == "true" || v == "t" || v == "1" {
		return true, nil
	}

	if v == "false" || v == "f" || v == "0" {
		return false, nil
	}

	return false, errors.New("invalid boolean value")
}

func apiURLsGet(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	c, err := parseCursor(r.FormValue("cursor"))
	if err != nil {
		writeJSONError(w, "invalid cursor value", http.StatusBadRequest)
		return
	}

	lim, err := parseInt(r.FormValue("limit"), 100)
	if err != nil || lim <= 0 || lim > 10000 {
		writeJSONError(w, "invalid limit value", http.StatusBadRequest)
		return
	}

	ig, err := parseBool(r.FormValue("include-generated-names"), false)
	if err != nil {
		writeJSONError(w, "invalid include-generated-names value", http.StatusBadRequest)
		return
	}

	res := msgRoutes{
		Ok: true,
	}

	iter := ctx.List(c)
	defer iter.Release()

	for iter.Next() {
		// if we should be ignoring generated links, skip over that range.
		if !ig && isGenerated(iter.Name()) {
			iter.Seek(postGenCursor)
			if !iter.Valid() {
				break
			}
		}

		res.Routes = append(res.Routes, &routeWithName{
			Name:  iter.Name(),
			Route: iter.Route(),
		})

		if len(res.Routes) == lim {
			break
		}
	}

	if iter.Next() {
		res.Next = base64.URLEncoding.EncodeToString([]byte(iter.Name()))
	}

	if err := iter.Error(); err != nil {
		writeJSONBackendError(w, err)
		return
	}

	writeJSON(w, &res, http.StatusOK)
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
		writeJSONError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusOK) // fix
	}
}

func apiURLs(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		apiURLsGet(ctx, w, r)
	default:
		writeJSONError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusOK) // fix
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
