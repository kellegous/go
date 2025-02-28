package backend

import (
	"context"

	"github.com/kellegous/go/internal"
)

type Backend interface {
	Close() error
	Get(ctx context.Context, id string) (*internal.Route, error)
	Put(ctx context.Context, key string, route *internal.Route) error
	Del(ctx context.Context, id string) error
	GetAll(ctx context.Context) (map[string]internal.Route, error)
	List(ctx context.Context, start string) (internal.RouteIterator, error)
	NextID(ctx context.Context) (uint64, error)
}
