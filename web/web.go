package web

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/kellegous/go/context"
	"github.com/syndtr/goleveldb/leveldb"
)

func makeName() string {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, rand.Int63())
	return base64.URLEncoding.EncodeToString(buf.Bytes())
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
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Panic(err)
	}
}

func writeJSONError(w http.ResponseWriter, error string, status int) {
	writeJSON(w, map[string]interface{}{
		"error": error,
	}, status)
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

func (h *apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := parseName("/api/url/", r.URL.Path)
	if p == "" {
		writeJSONError(w,
			http.StatusText(http.StatusNotFound),
			http.StatusNotFound)
		return
	}

	if r.Method == "POST" {
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

		rt := context.Route{
			URL:  req.URL,
			Time: time.Now(),
		}

		if err := h.ctx.Put(p, &rt); err != nil {
			log.Panic(err)
		}

		writeJSONRoute(w, p, &rt)
	} else if r.Method == "GET" {
		rt, err := h.ctx.Get(p)
		if err == leveldb.ErrNotFound {
			writeJSONError(w, "no such route", http.StatusNotFound)
			return
		} else if err != nil {
			log.Panic(err)
		}

		writeJSONRoute(w, p, rt)
	} else {
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
		http.Redirect(w, r,
			fmt.Sprintf("/edit/%s", makeName()),
			http.StatusTemporaryRedirect)
		return
	}

	rt, err := h.ctx.Get(p)
	if err == leveldb.ErrNotFound {
		http.Redirect(w, r,
			fmt.Sprintf("/edit/%s", p),
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
	p := parseName("/edit/", r.URL.Path)

	if p == "" {
		http.Redirect(w, r,
			fmt.Sprintf("/edit/%s", makeName()),
			http.StatusTemporaryRedirect)
		return
	}

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
