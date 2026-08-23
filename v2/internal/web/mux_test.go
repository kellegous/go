package web

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestURIToKeyword(t *testing.T) {
	tests := []struct {
		Name     string
		URI      string
		Expected string
	}{
		{Name: "empty", URI: "", Expected: ""},
		{Name: "root", URI: "/", Expected: ""},
		{Name: "keyword", URI: "/Example/path", Expected: "example"},
		{Name: "leading slash only", URI: "//Example/path", Expected: ""},
		{Name: "no leading slash", URI: "Example/path", Expected: "example"},
		{Name: "trailing slash", URI: "/Example/", Expected: "example"},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if got := uriToKeyword(tt.URI); got != tt.Expected {
				t.Errorf("uriToKeyword(%q) = %q, want %q", tt.URI, got, tt.Expected)
			}
		})
	}
}

func TestMuxReservations(t *testing.T) {
	useHandle := func(m *Mux, pattern string) {
		m.Handle(pattern, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	}

	useHandleFunc := func(m *Mux, pattern string) {
		m.HandleFunc(pattern, func(http.ResponseWriter, *http.Request) {})
	}

	tests := []struct {
		Name     string
		Setup    func(*Mux)
		Expected []string
	}{
		{
			Name: "Handle",
			Setup: func(m *Mux) {
				useHandle(m, "/")
				useHandle(m, "/a/url/")
				useHandle(m, "/A/url/")
				useHandle(m, "/b")
				useHandle(m, "/B/")
			},
			Expected: []string{"a", "b"},
		},
		{
			Name: "HandleFunc",
			Setup: func(m *Mux) {
				useHandleFunc(m, "/")
				useHandleFunc(m, "/a/url/")
				useHandleFunc(m, "/A/url/")
				useHandleFunc(m, "/b")
				useHandleFunc(m, "/B/")
			},
			Expected: []string{"a", "b"},
		},
		{
			Name: "Mixed",
			Setup: func(m *Mux) {
				useHandle(m, "/a")
				useHandleFunc(m, "/a/url/")
				useHandleFunc(m, "/b")
				useHandle(m, "/B/url")
			},
			Expected: []string{"a", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			m := NewMux()
			tt.Setup(m)

			if got := m.Reserved(); !reflect.DeepEqual(got, tt.Expected) {
				t.Errorf("Reserved() = %v, expected %v", got, tt.Expected)
			}

			for _, keyword := range tt.Expected {
				if got := m.IsReserved(keyword); !got {
					t.Errorf("IsReserved(%q) = false, expected true", keyword)
				}
			}
		})
	}
}

func TestMuxServeHTTP(t *testing.T) {
	m := NewMux()
	m.HandleFunc("/exact", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("exact"))
	})
	m.Handle("/prefix/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("prefix"))
	}))

	tests := []struct {
		Name     string
		Path     string
		Expected string
	}{
		{Name: "exact route", Path: "/exact", Expected: "exact"},
		{Name: "prefix route", Path: "/prefix/child", Expected: "prefix"},
		{Name: "not found", Path: "/missing", Expected: "404 page not found\n"},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.Path, nil)
			res := httptest.NewRecorder()
			m.ServeHTTP(res, req)

			if got := res.Body.String(); got != tt.Expected {
				t.Errorf("response body = %q, want %q", got, tt.Expected)
			}
		})
	}
}
