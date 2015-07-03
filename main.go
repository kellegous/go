package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/syndtr/goleveldb/leveldb"

	"github.com/kellegous/go/context"
)

func MakeName() string {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, rand.Int63())
	return base64.URLEncoding.EncodeToString(buf.Bytes())
}

func ParseName(base, path string) string {
	t := path[len(base):]
	ix := strings.Index(t, "/")
	if ix == -1 {
		return t
	} else {
		return t[:ix]
	}
}

func WriteJson(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Panic(err)
	}
}

func WriteJsonError(w http.ResponseWriter, error string, status int) {
	WriteJson(w, map[string]interface{}{
		"error": error,
	}, status)
}

func WriteJsonRoute(w http.ResponseWriter, name string, rt *context.Route) {
	res := struct {
		Name string    `json:"name"`
		URL  string    `json:"url"`
		Time time.Time `json:"time"`
	}{
		name,
		rt.URL,
		rt.Time,
	}

	WriteJson(w, &res, http.StatusOK)
}

func ServeAsset(w http.ResponseWriter, r *http.Request, name string) {
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

type DefaultHandler struct {
	ctx *context.Context
}

func (h *DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := ParseName("/", r.URL.Path)

	if p == "" {
		http.Redirect(w, r,
			fmt.Sprintf("/edit/%s", MakeName()),
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

type EditHandler struct {
	ctx *context.Context
}

func (h *EditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := ParseName("/edit/", r.URL.Path)

	if p == "" {
		http.Redirect(w, r,
			fmt.Sprintf("/edit/%s", MakeName()),
			http.StatusTemporaryRedirect)
		return
	}

	ServeAsset(w, r, "index.html")
}

type ApiHandler struct {
	ctx *context.Context
}

func (h *ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := ParseName("/api/url/", r.URL.Path)
	if p == "" {
		WriteJsonError(w,
			http.StatusText(http.StatusNotFound),
			http.StatusNotFound)
		return
	}

	if r.Method == "POST" {
		var req struct {
			URL string `json:"url"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteJsonError(w, "invalid json", http.StatusBadRequest)
			return
		}

		if req.URL == "" {
			WriteJsonError(w, "url required", http.StatusBadRequest)
			return
		}

		rt := context.Route{
			URL:  req.URL,
			Time: time.Now(),
		}

		if err := h.ctx.Put(p, &rt); err != nil {
			log.Panic(err)
		}

		WriteJsonRoute(w, p, &rt)
	} else if r.Method == "GET" {
		rt, err := h.ctx.Get(p)
		if err == leveldb.ErrNotFound {
			WriteJsonError(w, "no such route", http.StatusNotFound)
			return
		} else if err != nil {
			log.Panic(err)
		}

		WriteJsonRoute(w, p, rt)
	} else {
		WriteJsonError(w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed)
	}
}

func main() {
	flagData := flag.String("data", "data", "data")
	flagAddr := flag.String("addr", ":8067", "addr")
	flag.Parse()

	ctx, err := context.Open(*flagData)
	if err != nil {
		log.Panic(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", &DefaultHandler{ctx})
	mux.Handle("/edit/", &EditHandler{ctx})
	mux.Handle("/api/url/", &ApiHandler{ctx})
	mux.HandleFunc("/s/", func(w http.ResponseWriter, r *http.Request) {
		ServeAsset(w, r, r.URL.Path[len("/s/"):])
	})

	log.Panic(http.ListenAndServe(*flagAddr, mux))
}
