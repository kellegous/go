package store

import (
	"context"
	"errors"
	"iter"

	"github.com/kellegous/golinks"
)

var ErrLinkNotfound = errors.New("link not found")

type Store interface {
	Close() error
	Get(ctx context.Context, prefix string) (*golinks.Link, error)
	Put(ctx context.Context, link *golinks.Link) error
	Delete(ctx context.Context, prefix string) error
	List(ctx context.Context, start string) iter.Seq2[*golinks.Link, error]
}
