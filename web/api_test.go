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

func (e *env) getAPI(m *msg, name string) error {
	return e.callAPI(m, "GET", name, nil)
}

func (e *env) postAPI(m *msg, name, url string) error {
	r := struct {
		URL string `json:"url"`
	}{
		url,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&r); err != nil {
		return err
	}

	return e.callAPI(m, "POST", name, &buf)
}

func (e *env) callAPI(m *msg, method, name string, body io.Reader) error {
	req, err := http.NewRequest(method, fmt.Sprintf("/api/url/%s", name), body)
	if err != nil {
		return err
	}

	res := mockResponse{
		header: map[string][]string{},
	}

	e.mux.ServeHTTP(&res, req)

	if err := json.NewDecoder(&res).Decode(&m); err != nil {
		return err
	}

	return nil
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

	mux := http.NewServeMux()

	Setup(mux, ctx)

	return &env{
		mux: mux,
		dir: dir,
		ctx: ctx,
	}, nil
}

func needEnv(t *testing.T) *env {
	e, err := newEnv()
	if err != nil {
		t.Fatal(err)
	}
	return e
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

func assertOkWithRoute(t *testing.T, m *msg, url string) {
	if !m.Ok {
		t.Fatal("expected OK message, but it's not OK")
	}

	if m.Error != "" {
		t.Fatalf("expected no error, but got %s", m.Error)
	}

	if m.Route == nil {
		t.Fatalf("Route is nil, expected one with url of %s", url)
	}

	if m.Route.URL != url {
		t.Fatalf("Expected url of %s, got %s", url, m.Route.URL)
	}
}

func assertOkWithNamedRoute(t *testing.T, m *msg, name, url string) {
	assertOkWithRoute(t, m, url)
	if m.Route.Name != name {
		t.Fatalf("expected name %s, got %s", name, m.Route.Name)
	}
}

func TestAPIGetNotFound(t *testing.T) {
	e := needEnv(t)
	defer e.destroy()

	var m msg
	names := []string{"", "nothing", "nothing/there"}
	for _, name := range names {
		if err := e.getAPI(&m, name); err != nil {
			t.Fatal(err)
		}
		assertJustOk(t, &m)
	}
}

func TestAPIPutThenGet(t *testing.T) {
	e := needEnv(t)
	defer e.destroy()

	var pm msg
	if err := e.postAPI(&pm, "xxx", "http://ex.com/"); err != nil {
		t.Fatal(err)
	}
	assertOkWithRoute(t, &pm, "http://ex.com/")

	var gm msg
	if err := e.getAPI(&gm, "xxx"); err != nil {
		t.Fatal(err)
	}
	assertOkWithNamedRoute(t, &gm, "xxx", "http://ex.com/")
}

func TestAPIDel(t *testing.T) {
	e := needEnv(t)
	defer e.destroy()

	var am msg
	if err := e.postAPI(&am, "yyy", ""); err != nil {
		t.Fatal(err)
	}
	assertJustOk(t, &am)

	var bm msg
	if err := e.postAPI(&bm, "yyy", "https://a.com/"); err != nil {
		t.Fatal(err)
	}
	assertOkWithNamedRoute(t, &bm, "yyy", "https://a.com/")

	var cm msg
	if err := e.postAPI(&cm, "yyy", ""); err != nil {
		t.Fatal(err)
	}
	assertJustOk(t, &cm)

	var dm msg
	if err := e.getAPI(&dm, "yyy"); err != nil {
		t.Fatal(err)
	}
	assertJustOk(t, &dm)
}

func TestAPIPutThenGetAuto(t *testing.T) {
	e := needEnv(t)
	defer e.destroy()

	var am msg
	if err := e.postAPI(&am, "", "http://b.com/"); err != nil {
		t.Fatal(err)
	}
	assertOkWithRoute(t, &am, "http://b.com/")

	var bm msg
	if err := e.getAPI(&bm, am.Route.Name); err != nil {
		t.Fatal(err)
	}
	assertOkWithNamedRoute(t, &bm, am.Route.Name, "http://b.com/")
}
