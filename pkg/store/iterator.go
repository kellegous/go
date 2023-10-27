package store

import "context"

type Iterator[T any] interface {
	Next(ctx context.Context) (T, error)
	Close() error
}
