package store

import "context"

type RouteIterator interface {
	Next(ctx context.Context) (*Route, error)
}
