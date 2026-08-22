package internal

import (
	"context"
	"encoding/json"
	"errors"
	"iter"
	"regexp"
	"testing"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/kellegous/golinks"
	"github.com/kellegous/golinks/internal/store"
	"github.com/kellegous/poop"
)

var marshalOptions = protojson.MarshalOptions{
	UseProtoNames: true,
	Indent:        "  ",
}

func describe[T proto.Message](t *testing.T, m T) string {
	b, err := marshalOptions.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func describeSlice[T proto.Message](t *testing.T, s []T) string {
	items := make([]json.RawMessage, 0, len(s))
	for _, item := range s {
		b, err := marshalOptions.Marshal(item)
		if err != nil {
			t.Fatal(err)
		}
		items = append(items, json.RawMessage(b))
	}

	b, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func allAreSame[T proto.Message](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !proto.Equal(a[i], b[i]) {
			return false
		}
	}
	return true
}

func collectWithErr[T any](it iter.Seq2[T, error]) ([]T, error) {
	var results []T
	for item, err := range it {
		if err != nil {
			return nil, poop.Chain(err)
		}
		results = append(results, item)
	}
	return results, nil
}

type TestAdapter interface {
	NewStore(t *testing.T) (store.Store, func() error)
}

func Test(t *testing.T, ta TestAdapter) {
	tests := []struct {
		name string
		fn   func(t *testing.T, ta TestAdapter)
	}{
		{"TestGetPut", testGetPut},
		{"TestList", testList},
		{"TestDelete", testDelete},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fn(t, ta)
		})
	}
}

func testGetPut(t *testing.T, ta TestAdapter) {
	s, cleanup := ta.NewStore(t)

	defer cleanup()

	a := &golinks.Link{
		Prefix:    "a",
		CreatedAt: timestamppb.New(time.Unix(0x420, 0)),
		Matches: []*golinks.Match{
			{
				Pattern: regexp.MustCompile("a/(.*)").String(),
				Url:     "https://a.com/a/$1",
			},
		},
	}

	b := &golinks.Link{
		Prefix:    "b",
		CreatedAt: timestamppb.New(time.Unix(0x666, 0)),
		Matches: []*golinks.Match{
			{
				Pattern: regexp.MustCompile("b/(.*)").String(),
				Url:     "https://b.com/b/$1",
			},
		},
	}

	if err := s.Put(t.Context(), a); err != nil {
		t.Fatal(err)
	}

	if err := s.Put(t.Context(), b); err != nil {
		t.Fatal(err)
	}

	if ac, err := s.Get(t.Context(), "a"); err != nil {
		t.Fatal(err)
	} else if !proto.Equal(a, ac) {
		t.Fatalf("expected %s, got %s", describe(t, a), describe(t, ac))
	}

	if bc, err := s.Get(t.Context(), "b"); err != nil {
		t.Fatal(err)
	} else if !proto.Equal(b, bc) {
		t.Fatalf("expected %s, got %s", describe(t, b), describe(t, bc))
	}

	if c, err := s.Get(t.Context(), "c"); !errors.Is(err, store.ErrLinkNotfound) || c != nil {
		t.Fatalf("expected %v, got %v", store.ErrLinkNotfound, err)
	}
}

func testList(t *testing.T, ta TestAdapter) {
	s, cleanup := ta.NewStore(t)
	defer cleanup()

	a := &golinks.Link{
		Prefix:    "a",
		CreatedAt: timestamppb.New(time.Unix(0x420, 0)),
		Matches: []*golinks.Match{
			{
				Pattern: regexp.MustCompile("a/(.*)").String(),
				Url:     "https://a.com/a/$1",
			},
		},
	}

	b := &golinks.Link{
		Prefix:    "b",
		CreatedAt: timestamppb.New(time.Unix(0x666, 0)),
		Matches: []*golinks.Match{
			{
				Pattern: regexp.MustCompile("b/(.*)").String(),
				Url:     "https://b.com/b/$1",
			},
		},
	}

	c := &golinks.Link{
		Prefix:    "c",
		CreatedAt: timestamppb.New(time.Unix(0x420420, 0)),
		Matches: []*golinks.Match{
			{
				Pattern: regexp.MustCompile("c/(.*)").String(),
				Url:     "https://c.com/c/$1",
			},
		},
	}

	ctx := t.Context()

	{ // Test that an empty store returns an empty iterator.
		items, err := collectWithErr(s.List(ctx, ""))
		if err != nil {
			t.Fatal(err)
		}

		if !allAreSame(items, []*golinks.Link{}) {
			t.Fatalf(
				"expected %s, got %s",
				describeSlice(t, []*golinks.Link{}),
				describeSlice(t, items),
			)
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

	{ // Test full iteration.
		contents, err := collectWithErr(s.List(ctx, ""))
		if err != nil {
			t.Fatal(err)
		}

		expected := []*golinks.Link{
			a,
			b,
			c,
		}

		if !allAreSame(contents, expected) {
			t.Fatalf(
				"expected %s got %s",
				describeSlice(t, expected),
				describeSlice(t, contents))
		}
	}

	{ // Test iteration from a start key. The start key is inclusive.
		contents, err := collectWithErr(s.List(ctx, "a"))
		if err != nil {
			t.Fatal(err)
		}

		expected := []*golinks.Link{
			a,
			b,
			c,
		}

		if !allAreSame(contents, expected) {
			t.Fatalf(
				"expected %s got %s",
				describeSlice(t, expected),
				describeSlice(t, contents))
		}
	}

	{ // Test iteration from the last start key.
		contents, err := collectWithErr(s.List(ctx, "c"))
		if err != nil {
			t.Fatal(err)
		}

		expected := []*golinks.Link{c}

		if !allAreSame(contents, expected) {
			t.Fatalf(
				"expected %s got %s",
				describeSlice(t, expected),
				describeSlice(t, contents))
		}
	}
}

func testDelete(t *testing.T, ta TestAdapter) {
	s, cleanup := ta.NewStore(t)
	defer cleanup()

	a := &golinks.Link{
		Prefix:    "a",
		CreatedAt: timestamppb.New(time.Unix(0x420, 0)),
		Matches: []*golinks.Match{
			{
				Pattern: regexp.MustCompile("a/(.*)").String(),
				Url:     "https://a.com/a/$1",
			},
		},
	}

	b := &golinks.Link{
		Prefix:    "b",
		CreatedAt: timestamppb.New(time.Unix(0x666, 0)),
		Matches: []*golinks.Match{
			{
				Pattern: regexp.MustCompile("b/(.*)").String(),
				Url:     "https://b.com/b/$1",
			},
		},
	}

	c := &golinks.Link{
		Prefix:    "c",
		CreatedAt: timestamppb.New(time.Unix(0x420420, 0)),
		Matches: []*golinks.Match{
			{
				Pattern: regexp.MustCompile("c/(.*)").String(),
				Url:     "https://c.com/c/$1",
			},
		},
	}

	ctx := context.Background()

	links := []*golinks.Link{a, b, c}

	for i, r := range links {
		if err := s.Del(ctx, r.Prefix); !errors.Is(err, store.ErrLinkNotfound) {
			t.Fatalf("(%d) expected %v, got %v",
				i,
				store.ErrLinkNotfound,
				err)
		}
	}

	for _, r := range links {
		if err := s.Put(ctx, r); err != nil {
			t.Fatal(err)
		}
	}

	for _, r := range links {
		if err := s.Del(ctx, r.Prefix); err != nil {
			t.Fatal(err)
		}
	}

	for _, r := range links {
		if _, err := s.Get(ctx, r.Prefix); !errors.Is(err, store.ErrLinkNotfound) {
			t.Fatalf("expected %v, got %v", store.ErrLinkNotfound, err)
		}
	}
}
