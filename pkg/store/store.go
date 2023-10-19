package store

import (
	"context"

	"github.com/kellegous/golinks/pkg/internal"
)

type Store interface {
	Close() error
	Get(ctx context.Context) (*internal.Route, error)
}
