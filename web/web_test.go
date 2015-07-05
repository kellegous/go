package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/kellegous/go/context"
)

type env struct {
	mux *http.ServeMux
	dir string
	ctx *context.Context
}

func (e *env) destroy() {
	os.RemoveAll(e.dir)
}

func (e *env) callAPI(method, name string, body io.Reader) (*msg, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("/api/url/%s", name), body)
	if err != nil {
		return nil, err
	}

	res := mockResponse{
		header: map[string][]string{},
	}

	e.mux.ServeHTTP(&res, req)

	var m msg

	if err := json.NewDecoder(&res).Decode(&m); err != nil {
		return nil, err
	}

	return &m, nil
}

func newEnv() (*env, error) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}

	ctx, err := context.Open(filepath.Join(dir, "data"))
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}

	return &env{
		mux: allRoutes(ctx),
		dir: dir,
		ctx: ctx,
	}, nil
}

type mockResponse struct {
	header http.Header
	bytes.Buffer
	status int
}

func (r *mockResponse) Header() http.Header {
	return r.header
}

func (r *mockResponse) WriteHeader(status int) {
	r.status = status
}

func assertJustOk(t *testing.T, m *msg) {
	if !m.Ok {
		t.Fatal("expected OK message, but it's not OK")
	}

	if m.Error != "" {
		t.Fatalf("expected no error, but got %s", m.Error)
	}

	if m.Route != nil {
		t.Fatalf("expected no route, got %v", m.Route)
	}
}

func TestAPIGet(t *testing.T) {
	e, err := newEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer e.destroy()

	names := []string{"", "nothing", "nothing/there"}
	for _, name := range names {
		m, err := e.callAPI("GET", name, nil)
		if err != nil {
			t.Fatal(err)
		}
		assertJustOk(t, m)
	}
}
