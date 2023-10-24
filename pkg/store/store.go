package store

import (
	"context"
)

type Store interface {
	Close() error
	GetForURI(ctx context.Context, uri string) ([]*Route, error)
	Get(ctx context.Context, pattern string) (*Route, error)
	Put(ctx context.Context) error
	List(ctx context.Context, start string) (RouteIterator, error)
	Search(ctx context.Context, terms []string) ([]*Route, error)
}
