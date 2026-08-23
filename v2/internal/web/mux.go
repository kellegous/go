package web

import (
	"net/http"
	"strings"
)

type Mux struct {
	mux *http.ServeMux
	// NOTE: registration while serving is not supported.
	reserved map[string]bool
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

func (m *Mux) Handle(pattern string, handler http.Handler) {
	m.mux.Handle(pattern, handler)
	m.reserved[uriToKeyword(pattern)] = true
}

func (m *Mux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	m.mux.HandleFunc(pattern, handler)
	m.reserved[uriToKeyword(pattern)] = true
}

func (m *Mux) IsReserved(keyword string) bool {
	return m.reserved[strings.ToLower(keyword)]
}

func NewMux() *Mux {
	return &Mux{
		mux: http.NewServeMux(),
	}
}

func uriToKeyword(uri string) string {
	uri = strings.TrimPrefix(uri, "/")
	kw, _, _ := strings.Cut(uri, "/")
	return strings.ToLower(kw)
}
