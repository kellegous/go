package store

import (
	"context"
	"errors"
	"regexp"
)

var ErrorRouteNotfound = errors.New("route not found")

type Store interface {
	Close() error
	GetForURI(ctx context.Context, uri string) ([]*Route, error)
	Get(ctx context.Context, pattern *regexp.Regexp) (*Route, error)
	Put(ctx context.Context, r *Route) error
	List(ctx context.Context, start string) (RouteIterator, error)
	Search(ctx context.Context, terms []string) ([]*Route, error)
}
