package leveldb

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stgarf/go-links/internal"
)

func TestGetPut(t *testing.T) {
	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	backend, err := New(filepath.Join(tmp, "data"))
	if err != nil {
		t.Fatal(err)
	}
	defer backend.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if _, err := backend.Get(ctx, "not_found"); err != internal.ErrRouteNotFound {
		t.Fatalf("expected ErrRouteNotFound, got \"%v\"", err)
	}

	a := &internal.Route{
		URL:  "http://www.kellegous.com/",
		Time: time.Now(),
	}

	if err := backend.Put(ctx, "key", a); err != nil {
		t.Fatal(err)
	}

	b, err := backend.Get(ctx, "key")
	if err != nil {
		t.Fatal(err)
	}

	if b.URL != a.URL {
		t.Fatalf("expected URL of %s, got %s", a.URL, b.URL)
	}

	if !b.Time.Equal(a.Time) {
		t.Fatalf("expected Time of %s, got %s", a.Time, b.Time)
	}
}

func TestNextID(t *testing.T) {
	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	backend, err := New(filepath.Join(tmp, "data"))
	if err != nil {
		t.Fatal(err)
	}
	defer backend.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var e uint64 = 1
	for i := 0; i < 501; i++ {
		r, err := backend.NextID(ctx)
		if err != nil {
			t.Fatal(err)
		}

		if r != e {
			t.Fatalf("expected %d, got %d", e, r)
		}

		e++
	}
}

func TestEmptyList(t *testing.T) {
	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	backend, err := New(filepath.Join(tmp, "data"))
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	it, err := backend.List(ctx, "")
	defer it.Release()

	if it.Valid() {
		t.Fatal("Expected iterator to be invalid at start")
	}

	if it.Next() {
		t.Fatal("Expected there to be no next")
	}

	if err := it.Error(); err != nil {
		t.Fatal(err)
	}
}

func putRoutes(ctx context.Context, backend *Backend, names ...string) error {
	for _, name := range names {
		if err := backend.Put(ctx, name, &internal.Route{
			URL:  fmt.Sprintf("http://%s/", name),
			Time: time.Unix(0, 420),
		}); err != nil {
			return err
		}
	}
	return nil
}

func mustBeIterOf(t *testing.T, iter internal.RouteIterator, names ...string) {
	defer iter.Release()

	if iter.Valid() {
		t.Fatal("expected Iter to be invalid at start")
	}

	if err := iter.Error(); err != nil {
		t.Fatal("expected Iter not to begin with error")
	}

	if iter.Name() != "" {
		t.Fatalf("expected empty name but got \"%s\"", iter.Name())
	}

	if iter.Route() != nil {
		t.Fatalf("expected empty route but got %v", iter.Route())
	}

	for i, name := range names {
		if !iter.Next() {
			t.Fatalf("at item %d, expected more items", i)
		}

		if !iter.Valid() {
			t.Fatalf("on item %d, expected a valid iterator", i)
		}

		if iter.Name() != name {
			t.Fatalf("expected name of %s, got %s", name, iter.Name())
		}

		if iter.Route().URL != fmt.Sprintf("http://%s/", name) {
			t.Fatalf("expected route to have URL of http://%s/ got %s", name, iter.Route().URL)
		}
	}

	if iter.Next() {
		t.Fatal("iterator has too many items")
	}

	if iter.Valid() {
		t.Fatal("iterator should not be valid at end")
	}
}

func TestList(t *testing.T) {
	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	backend, err := New(filepath.Join(tmp, "data"))
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := putRoutes(ctx, backend, "a", "c", "d"); err != nil {
		t.Fatal(err)
	}

	iter, err := backend.List(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	mustBeIterOf(t, iter, "a", "c", "d")

	iter, err = backend.List(ctx, "b")
	if err != nil {
		t.Fatal(err)
	}
	mustBeIterOf(t, iter, "c", "d")

	iter, err = backend.List(ctx, "z")
	if err != nil {
		t.Fatal(err)
	}
	mustBeIterOf(t, iter)
}
