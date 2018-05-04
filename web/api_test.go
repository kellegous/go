package web

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/HALtheWise/o-links/context"
)

type urlReq struct {
	URL string `json:"url"`
}

type env struct {
	mux *http.ServeMux
	ctx *context.Context
}

func (e *env) destroy(t *testing.T) {
	err := e.ctx.DropTable()
	if err != nil {
		t.Errorf("Unable to drop table: %v", err)
	}
	e.ctx.Close()
	if err != nil {
		t.Errorf("Unable to close database: %v", err)
	}
}

func (e *env) get(path string) (*mockResponse, error) {
	return e.call("GET", path, nil)
}

func (e *env) post(path string, body interface{}) (*mockResponse, error) {
	return e.callWithJSON("POST", path, body)
}

func (e *env) callWithJSON(method, path string, body interface{}) (*mockResponse, error) {
	var r io.Reader

	if body != nil {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
		r = &buf
	}

	return e.call(method, path, r)
}

func (e *env) call(method, path string, body io.Reader) (*mockResponse, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	res := &mockResponse{
		header: map[string][]string{},
	}

	e.mux.ServeHTTP(res, req)

	return res, nil
}

func newEnv() (*env, error) {
	ctx, err := context.OpenTestCtx()
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	Setup(mux, ctx)

	return &env{
		mux: mux,
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

func mustBeSameNamedRoute(t *testing.T, a, b *routeWithName) {
	if a.Name != b.Name || a.URL != b.URL {
		t.Errorf("routes are not same: %v vs %v", a, b)
	}
	// TODO: Also check creation, modified, and deleted times
}

func mustBeRouteOf(t *testing.T, rt *context.Route, url string) {
	if rt == nil {
		t.Fatal("route is nil")
	}

	if rt.URL != url {
		t.Fatalf("expected url of %s, got %s", url, rt.URL)
	}

	if rt.CreatedAt.IsZero() {
		t.Fatal("route time is empty")
	}
}

func mustBeNamedRouteOf(t *testing.T, rt *routeWithName, name, url string) {
	mustBeRouteOf(t, rt.Route, url)
	if rt.Name != name {
		t.Fatalf("expected name of %s, got %s", name, rt.Name)
	}
}

func mustBeOk(t *testing.T, ok bool) {
	if !ok {
		t.Fatal("response is not ok")
	}
}

func mustBeErr(t *testing.T, m *msgErr) {
	if m.Ok {
		t.Fatal("response is ok, should be err")
	}

	if m.Error == "" {
		t.Fatal("expected an Error, but it is empty")
	}
}

func mustHaveStatus(t *testing.T, res *mockResponse, status int) {
	if res.status != status {
		t.Fatalf("expected response status %d, got %d", status, res.status)
	}
}

func TestAPIGetNotFound(t *testing.T) {
	e := needEnv(t)
	defer e.destroy(t)

	names := map[string]int{
		"":              http.StatusBadRequest,
		"nothing":       http.StatusNotFound,
		"nothing/there": http.StatusNotFound,
	}

	for name, status := range names {
		res, err := e.get(fmt.Sprintf("/api/url/%s", name))
		if err != nil {
			t.Fatal(err)
		}

		mustHaveStatus(t, res, status)

		var m msgErr
		if err := json.NewDecoder(res).Decode(&m); err != nil {
			t.Fatal(err)
		}

		mustBeErr(t, &m)
	}
}

func TestAPIPutThenGet(t *testing.T) {
	e := needEnv(t)
	defer e.destroy(t)

	res, err := e.post("/api/url/xxx", &urlReq{
		URL: "http://ex.com/",
	})
	if err != nil {
		t.Fatal(err)
	}

	mustHaveStatus(t, res, http.StatusOK)

	var pm msgRoute
	if err := json.NewDecoder(res).Decode(&pm); err != nil {
		t.Fatal(err)
	}

	mustBeOk(t, pm.Ok)
	mustBeNamedRouteOf(t, pm.Route, "xxx", "http://ex.com/")

	res, err = e.get("/api/url/xxx")
	if err != nil {
		t.Fatal(err)
	}

	mustHaveStatus(t, res, http.StatusOK)

	var gm msgRoute
	if err := json.NewDecoder(res).Decode(&gm); err != nil {
		t.Fatal(err)
	}

	mustBeOk(t, gm.Ok)
	mustBeNamedRouteOf(t, pm.Route, "xxx", "http://ex.com/")
}

func TestBadPuts(t *testing.T) {
	e := needEnv(t)
	defer e.destroy(t)

	var m msgErr

	res, err := e.call("POST", "/api/url/yyy", bytes.NewBufferString("not json"))
	if err != nil {
		t.Fatal(err)
	}
	mustHaveStatus(t, res, http.StatusBadRequest)

	if err := json.NewDecoder(res).Decode(&m); err != nil {
		t.Fatal(err)
	}
	mustBeErr(t, &m)

	res, err = e.post("/api/url/yyy", &urlReq{})
	if err != nil {
		t.Fatal(err)
	}
	mustHaveStatus(t, res, http.StatusBadRequest)

	if err := json.NewDecoder(res).Decode(&m); err != nil {
		t.Fatal(err)
	}
	mustBeErr(t, &m)

	res, err = e.post("/api/url/yyy", &urlReq{"not a URL"})
	if err != nil {
		t.Fatal(err)
	}
	mustHaveStatus(t, res, http.StatusBadRequest)

	if err := json.NewDecoder(res).Decode(&m); err != nil {
		t.Fatal(err)
	}
	mustBeErr(t, &m)
}

func TestAPIDel(t *testing.T) {
	e := needEnv(t)
	defer e.destroy(t)

	if err := e.ctx.Put("xxx", &context.Route{
		URL:       "http://ex.com/",
		CreatedAt: time.Now(),
	}); err != nil {
		t.Fatal(err)
	}

	res, err := e.call("DELETE", "/api/url/xxx", nil)
	if err != nil {
		t.Fatal(err)
	}

	mustHaveStatus(t, res, http.StatusOK)

	var m msg
	if err := json.NewDecoder(res).Decode(&m); err != nil {
		t.Fatal(err)
	}
	mustBeOk(t, m.Ok)

	if _, err := e.ctx.Get("xxx"); err != sql.ErrNoRows {
		t.Fatal("expected xxx to be deleted")
	}
}

func TestAPIPutThenGetAuto(t *testing.T) {
	e := needEnv(t)
	defer e.destroy(t)

	res, err := e.post("/api/url/", &urlReq{URL: "http://b.com/"})
	if err != nil {
		t.Fatal(err)
	}

	mustHaveStatus(t, res, http.StatusOK)

	var am msgRoute
	if err := json.NewDecoder(res).Decode(&am); err != nil {
		t.Fatal(err)
	}
	mustBeOk(t, am.Ok)
	mustBeRouteOf(t, am.Route.Route, "http://b.com/")

	res, err = e.get(fmt.Sprintf("/api/url/%s", am.Route.Name))
	if err != nil {
		t.Fatal(err)
	}

	mustHaveStatus(t, res, http.StatusOK)

	var bm msgRoute
	if err := json.NewDecoder(res).Decode(&bm); err != nil {
		t.Fatal(err)
	}
	mustBeOk(t, bm.Ok)
	mustBeNamedRouteOf(t, bm.Route, am.Route.Name, "http://b.com/")
}

func getLinksTest(e *env, params url.Values) ([]*routeWithName, error) {
	res, err := e.get("/api/urls/?" + params.Encode())
	if err != nil {
		return nil, err
	}

	if res.status != http.StatusOK {
		return nil, fmt.Errorf("HTTP status: %d", res.status)
	}

	var m msgRoutes
	if err := json.NewDecoder(res).Decode(&m); err != nil {
		return nil, err
	}

	if !m.Ok {
		return nil, errors.New("response is not ok")
	}

	return m.Routes, nil
}

type listTest struct {
	Params url.Values
	Pages  []*routeWithName
}

func TestAPIList(t *testing.T) {
	e := needEnv(t)
	defer e.destroy(t)

	rts := []*routeWithName{
		{
			Name: "0",
			Route: &context.Route{
				URL:       "http://0.com/",
				CreatedAt: time.Now(),
			},
		},

		{
			Name: "1",
			Route: &context.Route{
				URL:       "http://1.com/",
				CreatedAt: time.Now(),
			},
		},

		{
			Name: ":cat",
			Route: &context.Route{
				URL:       "http://cat.com/",
				CreatedAt: time.Now(),
				Generated: true,
			},
		},

		{
			Name: "_dog",
			Route: &context.Route{
				URL:       "http://dog.com/",
				CreatedAt: time.Now(),
				Generated: true,
			},
		},

		{
			Name: "a",
			Route: &context.Route{
				URL:       "http://a.com/",
				CreatedAt: time.Now(),
			},
		},

		{
			Name: "b",
			Route: &context.Route{
				URL:       "http://b.com/",
				CreatedAt: time.Now(),
			},
		},
	}

	for _, rt := range rts {
		rt.Uid = fmt.Sprint(rand.Uint64())
		if err := e.ctx.Put(rt.Name, rt.Route); err != nil {
			t.Fatal(err)
		}
	}

	tests := []*listTest{
		{
			Params: url.Values(map[string][]string{}),
			Pages: []*routeWithName{
				rts[0], rts[1], rts[4], rts[5]},
		},
		{
			Params: url.Values(map[string][]string{
				"include-generated-names": {"true"},
			}),
			Pages: rts,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Test with ?%s", test.Params.Encode()),
			func(t *testing.T) {
				links, err := getLinksTest(e, test.Params)
				if err != nil {
					t.Fatal(err)
				}

				expected := test.Pages

				if len(links) != len(expected) {
					t.Fatalf("length mismatch expected %d got %d", len(expected), len(links))
				}

				for j, m := 0, len(links); j < m; j++ {
					mustBeSameNamedRoute(t, links[j], expected[j])
				}

			})
	}
}

func TestBadList(t *testing.T) {
	e := needEnv(t)
	defer e.destroy(t)

	tests := map[string]int{
		url.Values{
			"include-generated-names": {"butter"},
		}.Encode(): http.StatusBadRequest,
	}

	for params, status := range tests {
		res, err := e.get("/api/urls/?" + params)
		if err != nil {
			t.Fatal(err)
		}

		mustHaveStatus(t, res, status)

		var m msgErr
		if err := json.NewDecoder(res).Decode(&m); err != nil {
			t.Fatal(err)
		}

		mustBeErr(t, &m)
	}
}
