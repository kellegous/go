package web

import (
	"context"
	"iter"
	"sync"
	"testing"

	"connectrpc.com/connect"

	"github.com/kellegous/golinks"
	"github.com/kellegous/golinks/internal/store"
)

type memoryStore struct {
	mu    sync.Mutex
	links map[string]*golinks.Link
}

var _ store.Store = (*memoryStore)(nil)

func newMemoryStore() *memoryStore {
	return &memoryStore{links: make(map[string]*golinks.Link)}
}

func (s *memoryStore) Close() error {
	return nil
}

func (s *memoryStore) Get(_ context.Context, prefix string) (*golinks.Link, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	link, ok := s.links[prefix]
	if !ok {
		return nil, store.ErrLinkNotfound
	}
	return link, nil
}

func (s *memoryStore) Put(_ context.Context, link *golinks.Link) (*golinks.Link, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.links[link.GetPrefix()] = link
	return link, nil
}

func (s *memoryStore) Delete(_ context.Context, prefix string) (*golinks.Link, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	link, ok := s.links[prefix]
	if !ok {
		return nil, store.ErrLinkNotfound
	}
	delete(s.links, prefix)
	return link, nil
}

func (s *memoryStore) List(_ context.Context, _ string) iter.Seq2[*golinks.Link, error] {
	return func(yield func(*golinks.Link, error) bool) {
		s.mu.Lock()
		links := make([]*golinks.Link, 0, len(s.links))
		for _, link := range s.links {
			links = append(links, link)
		}
		s.mu.Unlock()

		for _, link := range links {
			if !yield(link, nil) {
				return
			}
		}
	}
}

func TestServerPutAndGet(t *testing.T) {
	s := &server{store: newMemoryStore()}
	want := &golinks.Link{Prefix: "go", Matches: []*golinks.Match{{Pattern: "", Url: "https://example.com"}}}

	putRes, err := s.Put(context.Background(), connect.NewRequest(&golinks.PutReq{Link: want}))
	if err != nil {
		t.Fatalf("Put() error = %v", err)
	}
	if putRes.Msg.GetLink() != want {
		t.Fatalf("Put() link = %p, want %p", putRes.Msg.GetLink(), want)
	}

	getRes, err := s.Get(context.Background(), connect.NewRequest(&golinks.GetReq{Prefix: "go"}))
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if getRes.Msg.GetLink() != want {
		t.Fatalf("Get() link = %p, want %p", getRes.Msg.GetLink(), want)
	}
}

func TestServerDelete(t *testing.T) {
	s := &server{store: newMemoryStore()}
	want := &golinks.Link{Prefix: "go", Matches: []*golinks.Match{{Pattern: "", Url: "https://example.com"}}}
	if _, err := s.store.Put(context.Background(), want); err != nil {
		t.Fatal(err)
	}

	res, err := s.Delete(context.Background(), connect.NewRequest(&golinks.DeleteReq{Prefix: "go"}))
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if res.Msg.GetLink() != want {
		t.Fatalf("Delete() link = %p, want %p", res.Msg.GetLink(), want)
	}

	_, err = s.Get(context.Background(), connect.NewRequest(&golinks.GetReq{Prefix: "go"}))
	if connect.CodeOf(err) != connect.CodeNotFound {
		t.Fatalf("Get() after Delete() code = %v, want %v", connect.CodeOf(err), connect.CodeNotFound)
	}
}

func TestServerNotFound(t *testing.T) {
	s := &server{store: newMemoryStore()}

	tests := []struct {
		Name string
		Call func() error
	}{
		{
			Name: "Get",
			Call: func() error {
				_, err := s.Get(context.Background(), connect.NewRequest(&golinks.GetReq{Prefix: "missing"}))
				return err
			},
		},
		{
			Name: "Delete",
			Call: func() error {
				_, err := s.Delete(context.Background(), connect.NewRequest(&golinks.DeleteReq{Prefix: "missing"}))
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if code := connect.CodeOf(tt.Call()); code != connect.CodeNotFound {
				t.Fatalf("error code = %v, want %v", code, connect.CodeNotFound)
			}
		})
	}
}
