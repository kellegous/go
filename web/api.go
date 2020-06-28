package web

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/stgarf/go-links/backend"
	"github.com/stgarf/go-links/internal"
)

const (
	alpha = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	errInvalidURL        = errors.New("Invalid URL")
	errRedirectLoop      = errors.New("I'm sorry, Dave. I'm afraid I can't do that")
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

// Advance to the next id and encode it as an ID.
func nextEncodedID(ctx context.Context, backend backend.Backend) (string, error) {
	id, err := backend.NextID(ctx)
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
	case "http", "https", "mailto", "ftp", "slack", "ssh", "zoommtg", "zoomus":
		break
	default:
		log.Printf("Invalid scheme for URL %s", u)
		return errInvalidURL
	}

	if r.Host == u.Host {
		return errRedirectLoop
	}

	return nil
}

func apiURLPost(backend backend.Backend, host string, w http.ResponseWriter, r *http.Request) {
	log.Printf("POST API request for url %s from %s:", r.URL.Path, r.RemoteAddr)
	p, _ := parseName("/api/url/", r.URL.Path)

	var req struct {
		URL  string `json:"url"`
		Hits string `json:"hits"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "invalid json", http.StatusBadRequest)
		log.Debugf("Bad json, failed to decode request body: %+v", r.Body)
		return
	}

	if req.URL == "" {
		writeJSONError(w, "url required", http.StatusBadRequest)
		log.Debugf("Url required: %+v", r.URL.RequestURI())
		return
	}

	if isBannedName(p) {
		writeJSONError(w, "name cannot be used", http.StatusBadRequest)
		log.Debugf("Banned named attempted: %+v", r.URL.RequestURI())
		return
	}

	if err := validateURL(r, req.URL); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// If no name is specified, an ID must be generated.
	if p == "" {
		var err error
		p, err = nextEncodedID(ctx, backend)
		if err != nil {
			writeJSONBackendError(w, err)
			return
		}
	}

	rt := internal.Route{
		URL:  req.URL,
		Time: time.Now(),
		Hits: req.Hits,
	}

	if err := backend.Put(ctx, p, &rt); err != nil {
		writeJSONBackendError(w, err)
		return
	}

	writeJSONRoute(w, p, &rt, host)
}

func apiURLGet(backend backend.Backend, host string, w http.ResponseWriter, r *http.Request) {
	p, _ := parseName("/api/url/", r.URL.Path)
	log.Debugf("GET API request for url %s from %s:", r.URL.Path, r.RemoteAddr)

	if p == "" {
		writeJSONError(w, "no name given", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	rt, err := backend.Get(ctx, p)
	if errors.Is(err, internal.ErrRouteNotFound) {
		writeJSONError(w, "Not Found", http.StatusNotFound)
		return
	} else if err != nil {
		writeJSONBackendError(w, err)
		return
	}

	writeJSONRoute(w, p, rt, host)
}

func apiURLDelete(backend backend.Backend, w http.ResponseWriter, r *http.Request) {
	p, _ := parseName("/api/url/", r.URL.Path)

	if p == "" {
		writeJSONError(w, "name required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := backend.Del(ctx, p); err != nil {
		writeJSONBackendError(w, err)
		return
	}
	log.Printf("DELETE API request for url %s from %s:", r.URL.Path, r.RemoteAddr)
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

func apiURLsGet(backend backend.Backend, host string, w http.ResponseWriter, r *http.Request) {
	log.Printf("GET API request (ALL URLS) for url %s from %s", r.URL.Path, r.RemoteAddr)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	iter, err := backend.List(ctx, string(c))
	if err != nil {
		writeJSONBackendError(w, err)
		return
	}
	defer iter.Release()

	for iter.Next() {
		// if we should be ignoring generated links, skip over that range.
		if !ig && isGenerated(iter.Name()) {
			iter.Seek(string(postGenCursor))
			if !iter.Valid() {
				break
			}
		}

		r := routeWithName{
			Name:  iter.Name(),
			Route: iter.Route(),
		}

		if host != "" {
			r.SourceHost = host
		}

		res.Routes = append(res.Routes, &r)

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

func apiURL(backend backend.Backend, host string, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		log.Debug("Handling POST")
		apiURLPost(backend, host, w, r)
	case "GET":
		log.Debug("Handling GET")
		apiURLGet(backend, host, w, r)
	case "DELETE":
		log.Debug("Handling DELETE")
		apiURLDelete(backend, w, r)
	default:
		log.Warnf("Handling %s... to %s Strange. Fuzzer? Hacker?!", r.Method, r.URL.Path)
		writeJSONError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusOK) // fix
	}
}

func apiURLs(backend backend.Backend, host string, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Debug("Handling GET")
		apiURLsGet(backend, host, w, r)
	default:
		log.Warnf("Handling %s... to %s Strange. Fuzzer? Hacker?!", r.Method, r.URL.Path)
		writeJSONError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusOK) // fix
	}
}

// Setup ...
func Setup(m *http.ServeMux, backend backend.Backend, host string) {
	m.HandleFunc("/api/url/", func(w http.ResponseWriter, r *http.Request) {
		apiURL(backend, host, w, r)
	})

	m.HandleFunc("/api/urls/", func(w http.ResponseWriter, r *http.Request) {
		apiURLs(backend, host, w, r)
	})
}
