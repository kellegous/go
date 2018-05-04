package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
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
	errRedirectLoop = errors.New(" I'm sorry, Dave. I'm afraid I can't do that." +
		"Recursive links are not currently supported.")
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
	name := parseName("/api/url/", r.URL.Path)
	var req struct {
		URL           string `json:"url"`
		Uid           string `json:"uid"`
		Generated     bool   `json:"generated"`
		ModifiedCount int    `json:"modified_count"`
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

	// If no path is specified, a path must be generated.
	if name == "" {
		var err error
		name, err = generateLink(ctx, req.Uid)
		if err != nil {
			writeJSONBackendError(w, err)
			return
		}

		req.Generated = true
	}

	name, err := normalizeName(name)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
	}

	reqURL, err := validateURL(r, req.URL)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	rt := context.Route{
		URL:           reqURL,
		CreatedAt:     time.Now(),
		Uid:           req.Uid,
		Generated:     req.Generated,
		ModifiedCount: req.ModifiedCount,
	}

	// If a row with the name already exists, ctx.Get won't return an error. We then call ctx.Edit() instead.
	if _, err := ctx.GetUid(rt.Uid); err == nil {
		if err := ctx.Edit(&rt, name); err != nil {
			writeJSONBackendError(w, err)
			return
		}
	} else {
		if err := ctx.Put(name, &rt); err != nil {
			writeJSONBackendError(w, err)
			return
		}
	}
	writeJSONRoute(w, name, &rt)
}

func apiURLGet(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	name, err := normalizeName(parseName("/api/url/", r.URL.Path))
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	rt, err := ctx.Get(name)
	if err == sql.ErrNoRows {
		writeJSONError(w, "Not Found", http.StatusNotFound)
		return
	} else if err != nil {
		writeJSONBackendError(w, err)
		return
	}

	writeJSONRoute(w, name, rt)
}

func apiURLDelete(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	name, err := normalizeName(parseName("/api/url/", r.URL.Path))
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ctx.Del(name); err != nil {
		writeJSONBackendError(w, err)
		return
	}

	writeJSONOk(w)
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
	ig, err := parseBool(r.FormValue("include-generated-names"), false)
	if err != nil {
		writeJSONError(w, "invalid include-generated-names value", http.StatusBadRequest)
		return
	}

	res := msgRoutes{
		Ok: true,
	}

	links, err := ctx.GetAll()
	if err != nil {
		writeJSONBackendError(w, err)
	}

	sortedNames := make([]string, 0, len(links))
	for k := range links {
		sortedNames = append(sortedNames, k)
	}
	sort.Strings(sortedNames)

	for _, name := range sortedNames {
		// if we should be ignoring generated links, skip over that range.
		route := links[name]
		if !ig && route.Generated {
			continue
		}

		res.Routes = append(res.Routes, &routeWithName{
			Name:  name,
			Route: &route,
		})
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
