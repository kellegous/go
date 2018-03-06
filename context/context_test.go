package context

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

func TestGetPut(t *testing.T) {
	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	ctx, err := Open(filepath.Join(tmp, "data"))
	if err != nil {
		t.Fatal(err)
	}
	defer ctx.Close()

	if _, err := ctx.Get("not_found"); err != leveldb.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got \"%v\"", err)
	}

	a := &Route{
		URL:  "http://www.kellegous.com/",
		Time: time.Now(),
	}

	if err := ctx.Put("key", a); err != nil {
		t.Fatal(err)
	}

	b, err := ctx.Get("key")
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

func TestEmptyList(t *testing.T) {
	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	ctx, err := Open(filepath.Join(tmp, "data"))
	if err != nil {
		t.Fatal(err)
	}

	it := ctx.List(nil)
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

func putRoutes(ctx *Context, names ...string) error {
	for _, name := range names {
		if err := ctx.Put(name, &Route{
			URL:  fmt.Sprintf("http://%s/", name),
			Time: time.Unix(0, 420),
		}); err != nil {
			return err
		}
	}
	return nil
}

func mustBeIterOf(t *testing.T, iter *Iter, names ...string) {
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

	ctx, err := Open(filepath.Join(tmp, "data"))
	if err != nil {
		t.Fatal(err)
	}

	if err := putRoutes(ctx, "a", "c", "d"); err != nil {
		t.Fatal(err)
	}

	mustBeIterOf(t, ctx.List(nil), "a", "c", "d")
	mustBeIterOf(t, ctx.List([]byte{'b'}), "c", "d")
	mustBeIterOf(t, ctx.List([]byte{'z'}))
}
