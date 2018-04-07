package web

import (
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"database/sql"

	"github.com/HALtheWise/o-links/context"
	_ "github.com/lib/pq"
)

const (
	alpha = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	errInvalidURL   = errors.New("Invalid URL")
	errRedirectLoop = errors.New(" I'm sorry, Dave. I'm afraid I can't do that")
)

// Check that the given URL is suitable as a shortcut link.
func validateURL(r *http.Request, s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", errInvalidURL
	}

	if u.Host == "" && u.Path == "" {
		return "", errInvalidURL
	}

	switch u.Scheme {
	case "":
		u.Scheme = "http"
		break
	case "http", "https", "mailto", "ftp":
		break
	default:
		return "", errInvalidURL
	}

	if r.Host == u.Host {
		return "", errRedirectLoop
	}

	return u.String(), nil
}

/*What if someone wants to edit an existing row by name? We need to check with ctx.Get() and then edit with a ctx.Edit() method*/
func apiURLPost(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	p := parseName("/api/url/", r.URL.Path)

	var req struct {
		URL       string `json:"url"`
		Uid       string `json:"uid"`
		Generated bool   `json:"generated"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		writeJSONError(w, "url required", http.StatusBadRequest)
		return
	}

	if req.Uid == "" {
		req.Uid = fmt.Sprint(randsource.Uint32())
	}

	if isBannedName(p) {
		writeJSONError(w, "name cannot be used", http.StatusBadRequest)
		return
	}

	reqURL, err := validateURL(r, req.URL)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If no path is specified, a path must be generated.
	if p == "" {
		var err error
		p, err = generateLink(ctx, req.Uid)
		if err != nil {
			writeJSONBackendError(w, err)
			return
		}

		req.Generated = true
	}

	rt := context.Route{
		URL:       reqURL,
		CreatedAt: time.Now(),
		Uid:       req.Uid,
		Generated: req.Generated,
	}

	// If a row with the name already exists, ctx.Get won't return an error. We then call ctx.Edit() instead.
	if _, err := ctx.Get(p); err == nil {
		if err := ctx.Edit(p, rt.URL); err != nil {
			writeJSONBackendError(w, err)
			return
		}
	} else {
		if err := ctx.Put(p, &rt); err != nil {
			writeJSONBackendError(w, err)
			return
		}
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
	if err == sql.ErrNoRows {
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

	return base32.StdEncoding.DecodeString(v)
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

// Setup ...
func Setup(m *http.ServeMux, ctx *context.Context) {
	m.HandleFunc("/api/url/", func(w http.ResponseWriter, r *http.Request) {
		apiURL(ctx, w, r)
	})
}
