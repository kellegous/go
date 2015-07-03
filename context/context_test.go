package context

import (
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

func TestNextID(t *testing.T) {
	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	ctx, err := Open(filepath.Join(tmp, "data"))
	if err != nil {
		t.Fatal(err)
	}

	var e uint64 = 1
	for i := 0; i < 501; i++ {
		r, err := ctx.NextID()
		if err != nil {
			t.Fatal(err)
		}

		if r != e {
			t.Fatalf("expected %d, got %d", e, r)
		}

		e++
	}
}
