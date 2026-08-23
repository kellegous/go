package web

import (
	"net/http"
	"sort"
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
	if kw := uriToKeyword(pattern); kw != "" {
		m.reserved[kw] = true
	}
}

func (m *Mux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	m.mux.HandleFunc(pattern, handler)
	if kw := uriToKeyword(pattern); kw != "" {
		m.reserved[kw] = true
	}
}

func (m *Mux) IsReserved(keyword string) bool {
	return m.reserved[strings.ToLower(keyword)]
}

func (m *Mux) Reserved() []string {
	reserved := make([]string, 0, len(m.reserved))
	for keyword := range m.reserved {
		reserved = append(reserved, keyword)
	}
	sort.Strings(reserved)
	return reserved
}

func NewMux() *Mux {
	return &Mux{
		mux:      http.NewServeMux(),
		reserved: make(map[string]bool),
	}
}

func uriToKeyword(uri string) string {
	uri = strings.TrimPrefix(uri, "/")
	kw, _, _ := strings.Cut(uri, "/")
	return strings.ToLower(kw)
}
