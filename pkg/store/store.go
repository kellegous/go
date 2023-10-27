package store

import (
	"context"
	"errors"
	"regexp"
)

var (
	ErrRouteNotfound = errors.New("route not found")
	ErrIteratorDone  = errors.New("iterator done")
)

type Store interface {
	Close() error
	GetForPrefix(ctx context.Context, prefix string) (Iterator[*Route], error)
	Get(ctx context.Context, pattern *regexp.Regexp) (*Route, error)
	Put(ctx context.Context, r *Route) error
	List(ctx context.Context, opts ListOptions) (Iterator[Cursor[*Route]], error)
	// Search(ctx context.Context, terms []string) (Iterator[*Route], error)
}
