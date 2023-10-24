package leveldb

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/kellegous/golinks/pkg/store"
)

type testStore struct {
	dir string
	*Store
}

func createTestStore(t *testing.T) (*Store, func() error) {
	tmp, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}

	s, err := Open(tmp)
	if err != nil {
		t.Fatal(err)
	}

	return s, func() error {
		return os.RemoveAll(tmp)
	}
}

func routesAreSame(a, b *store.Route) bool {
	if a == b {
		return true
	} else if b == nil || a == nil {
		return false
	}
	return a.Pattern.String() == b.Pattern.String() &&
		a.URL == b.URL &&
		a.Time.Equal(b.Time)
}

func describeRoute(r *store.Route) []byte {
	b, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return b
}

func TestGetPut(t *testing.T) {
	s, cleanup := createTestStore(t)
	defer cleanup()

	ap := regexp.MustCompile("a/(.*)")
	bp := regexp.MustCompile("b/(.*)")

	a := &store.Route{
		Pattern: ap,
		URL:     "https://a.com/a/$1",
		Time:    time.Unix(420, 0),
	}

	b := &store.Route{
		Pattern: bp,
		URL:     "https://b.com/b/$1",
		Time:    time.Unix(69, 0),
	}

	if err := s.Put(
		context.Background(),
		a,
	); err != nil {
		t.Fatal(err)
	}

	if err := s.Put(
		context.Background(),
		b,
	); err != nil {
		t.Fatal(err)
	}

	if ac, err := s.Get(context.Background(), ap); err != nil {
		t.Fatal(err)
	} else if !routesAreSame(a, ac) {
		t.Fatalf("expected %s, got %s", describeRoute(a), describeRoute(ac))
	}

	if bc, err := s.Get(context.Background(), bp); err != nil {
		t.Fatal(err)
	} else if !routesAreSame(b, bc) {
		t.Fatalf("expected %s, got %s", describeRoute(b), describeRoute(bc))
	}

	if c, err := s.Get(context.Background(), regexp.MustCompile("c")); !errors.Is(err, store.ErrorRouteNotfound) || c != nil {
		t.Fatalf("expected %v, got %v", store.ErrorRouteNotfound, err)
	}
}
