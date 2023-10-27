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

func describe(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func iteratorToSlice[T any](
	ctx context.Context,
	it store.Iterator[T],
) ([]T, error) {
	var res []T
	for {
		t, err := it.Next(ctx)
		if errors.Is(err, store.ErrIteratorDone) {
			break
		} else if err != nil {
			return nil, err
		}

		res = append(res, t)
	}

	return res, nil
}

func allAreSame[T any](a, b []T, cmp func(a, b T) bool) bool {
	n := len(a)
	if n != len(b) {
		return false
	}
	for i := 0; i < n; i++ {
		if !cmp(a[i], b[i]) {
			return false
		}
	}
	return true
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
		t.Fatalf("expected %s, got %s", describe(a), describe(ac))
	}

	if bc, err := s.Get(context.Background(), bp); err != nil {
		t.Fatal(err)
	} else if !routesAreSame(b, bc) {
		t.Fatalf("expected %s, got %s", describe(b), describe(bc))
	}

	if c, err := s.Get(
		context.Background(),
		regexp.MustCompile("c"),
	); !errors.Is(err, store.ErrRouteNotfound) || c != nil {
		t.Fatalf("expected %v, got %v", store.ErrRouteNotfound, err)
	}
}

func TestList(t *testing.T) {
	s, cleanup := createTestStore(t)
	defer cleanup()

	a := &store.Route{
		Pattern: regexp.MustCompile("a/(.*)"),
		URL:     "https://a.com/a/$1",
		Time:    time.Unix(420, 0),
	}

	b := &store.Route{
		Pattern: regexp.MustCompile("b/(.*)"),
		URL:     "https://b.com/b/$1",
		Time:    time.Unix(69, 0),
	}

	c := &store.Route{
		Pattern: regexp.MustCompile("c/(.*)"),
		URL:     "https://c.com/c/$1",
		Time:    time.Unix(666, 0),
	}

	ctx := context.Background()

	// Test that an empty store returns an empty iterator.
	{
		it, err := s.List(ctx, store.ListOptions{})
		if err != nil {
			t.Fatal(err)
		}
		defer it.Close()

		contents, err := iteratorToSlice(ctx, it)
		if err != nil {
			t.Fatal(err)
		}

		expected := []store.Cursor[*store.Route]{}

		if !allAreSame(
			contents,
			expected,
			func(a, b store.Cursor[*store.Route]) bool {
				return routesAreSame(a.Route(), b.Route())
			}) {
			t.Fatalf(
				"expected %s got %s",
				describe(expected),
				describe(contents))
		}
	}

	if err := s.Put(ctx, c); err != nil {
		t.Fatal(err)
	}

	if err := s.Put(ctx, a); err != nil {
		t.Fatal(err)
	}

	if err := s.Put(ctx, b); err != nil {
		t.Fatal(err)
	}

	// Test full iteration.
	{
		it, err := s.List(ctx, store.ListOptions{})
		if err != nil {
			t.Fatal(err)
		}
		defer it.Close()

		contents, err := iteratorToSlice(ctx, it)
		if err != nil {
			t.Fatal(err)
		}

		expected := []store.Cursor[*store.Route]{&Cursor{a}, &Cursor{b}, &Cursor{c}}
		if !allAreSame(
			contents,
			expected,
			func(a, b store.Cursor[*store.Route]) bool {
				return routesAreSame(a.Route(), b.Route())
			}) {
			t.Fatalf(
				"expected %s got %s",
				describe(expected),
				describe(contents))
		}
	}

	// Test iteration from a cursor.
	{
		it, err := s.List(ctx, store.ListOptions{Cursor: encodeCursor(a.Pattern)})
		if err != nil {
			t.Fatal(err)
		}
		defer it.Close()

		contents, err := iteratorToSlice(ctx, it)
		if err != nil {
			t.Fatal(err)
		}

		expected := []store.Cursor[*store.Route]{
			&Cursor{b},
			&Cursor{c},
		}

		if !allAreSame(
			contents,
			expected,
			func(a, b store.Cursor[*store.Route]) bool {
				return routesAreSame(a.Route(), b.Route())
			}) {
			t.Fatalf(
				"expected %s got %s",
				describe(expected),
				describe(contents))
		}
	}

	// Test iteration from last cursor.
	{
		it, err := s.List(ctx, store.ListOptions{Cursor: encodeCursor(c.Pattern)})
		if err != nil {
			t.Fatal(err)
		}
		defer it.Close()

		contents, err := iteratorToSlice(ctx, it)
		if err != nil {
			t.Fatal(err)
		}

		expected := []store.Cursor[*store.Route]{}

		if !allAreSame(
			contents,
			expected,
			func(a, b store.Cursor[*store.Route]) bool {
				return routesAreSame(a.Route(), b.Route())
			}) {
			t.Fatalf(
				"expected %s got %s",
				describe(expected),
				describe(contents))
		}
	}
}

func TestGetForPrefix(t *testing.T) {
	s, cleanup := createTestStore(t)
	defer cleanup()

	a := &store.Route{
		Pattern: regexp.MustCompile("a/(.*)"),
		URL:     "https://a.com/a/$1",
		Time:    time.Unix(420, 0),
	}

	b := &store.Route{
		Pattern: regexp.MustCompile("b/1/(.*)"),
		URL:     "https://b.com/b/$1",
		Time:    time.Unix(69, 0),
	}

	c := &store.Route{
		Pattern: regexp.MustCompile("b/2/(.*)"),
		URL:     "https://b.com/b/$1",
		Time:    time.Unix(69, 0),
	}

	d := &store.Route{
		Pattern: regexp.MustCompile("c/(.*)"),
		URL:     "https://b.com/b/$1",
		Time:    time.Unix(69, 0),
	}

	ctx := context.Background()
	for _, route := range []*store.Route{a, b, c, d} {
		if err := s.Put(ctx, route); err != nil {
			t.Fatal(err)
		}
	}

	// Test iterating first prefix in store.
	{
		it, err := s.GetForPrefix(ctx, "a")
		if err != nil {
			t.Fatal(err)
		}
		defer it.Close()

		contents, err := iteratorToSlice(ctx, it)
		if err != nil {
			t.Fatal(err)
		}

		expected := []*store.Route{a}

		if !allAreSame(
			contents,
			expected,
			func(a, b *store.Route) bool {
				return routesAreSame(a, b)
			}) {
			t.Fatalf(
				"expected %s got %s",
				describe(expected),
				describe(contents))
		}
	}

	// Test iterating second prefix in store.
	{
		it, err := s.GetForPrefix(ctx, "b")
		if err != nil {
			t.Fatal(err)
		}
		defer it.Close()

		contents, err := iteratorToSlice(ctx, it)
		if err != nil {
			t.Fatal(err)
		}

		expected := []*store.Route{b, c}

		if !allAreSame(
			contents,
			expected,
			func(a, b *store.Route) bool {
				return routesAreSame(a, b)
			}) {
			t.Fatalf(
				"expected %s got %s",
				describe(expected),
				describe(contents))
		}
	}

	// Test iterating last prefix in store.
	{
		it, err := s.GetForPrefix(ctx, "c")
		if err != nil {
			t.Fatal(err)
		}
		defer it.Close()

		contents, err := iteratorToSlice(ctx, it)
		if err != nil {
			t.Fatal(err)
		}

		expected := []*store.Route{d}

		if !allAreSame(
			contents,
			expected,
			func(a, b *store.Route) bool {
				return routesAreSame(a, b)
			}) {
			t.Fatalf(
				"expected %s got %s",
				describe(expected),
				describe(contents))
		}
	}

	// Test iterating non-existent prefix in store.
	{
		it, err := s.GetForPrefix(ctx, "z")
		if err != nil {
			t.Fatal(err)
		}
		defer it.Close()

		contents, err := iteratorToSlice(ctx, it)
		if err != nil {
			t.Fatal(err)
		}

		expected := []*store.Route{}

		if !allAreSame(
			contents,
			expected,
			func(a, b *store.Route) bool {
				return routesAreSame(a, b)
			}) {
			t.Fatalf(
				"expected %s got %s",
				describe(expected),
				describe(contents))
		}
	}
}
