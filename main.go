package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	dbFilename = "keys.db"
)

type Route struct {
	Url  string
	Time time.Time
}

func (r *Route) Write(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, r.Time.UnixNano()); err != nil {
		return err
	}

	if _, err := w.Write([]byte(r.Url)); err != nil {
		return err
	}

	return nil
}

func (o *Route) Read(r io.Reader) error {
	var t int64
	if err := binary.Read(r, binary.LittleEndian, &t); err != nil {
		return err
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	o.Url = string(b)
	o.Time = time.Unix(0, t)
	return nil
}

type Context struct {
	path string
}

func (c *Context) Init() error {
	if _, err := os.Stat(c.path); err != nil {
		if err := os.MkdirAll(c.path, os.ModePerm); err != nil {
			return err
		}
	}

	db, err := openDb(c.path)
	if err != nil {
		return err
	}

	return db.Close()
}

func openDb(path string) (*leveldb.DB, error) {
	return leveldb.OpenFile(filepath.Join(path, dbFilename), nil)
}

func (c *Context) Get(key string) (*Route, error) {
	db, err := openDb(c.path)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	val, err := db.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}

	r := &Route{}
	if err := r.Read(bytes.NewBuffer(val)); err != nil {
		return nil, err
	}

	return r, nil
}

func (c *Context) Put(key string, r *Route) error {
	db, err := openDb(c.path)
	if err != nil {
		return err
	}
	defer db.Close()

	var buf bytes.Buffer

	if err := r.Write(&buf); err != nil {
		return err
	}

	return db.Put([]byte(key), buf.Bytes(), nil)
}

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

func WriteJsonRoute(w http.ResponseWriter, name string, rt *Route) {
	res := struct {
		Name string    `json:"name"`
		URL  string    `json:"url"`
		Time time.Time `json:"time"`
	}{
		name,
		rt.Url,
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
	ctx *Context
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
		rt.Url,
		http.StatusTemporaryRedirect)
}

type EditHandler struct {
	ctx *Context
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
	ctx *Context
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

		rt := Route{
			Url:  req.URL,
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

	ctx := &Context{
		path: *flagData,
	}

	if err := ctx.Init(); err != nil {
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
