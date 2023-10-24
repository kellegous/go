package leveldb

import (
	"context"
	"errors"
	"regexp"

	"github.com/syndtr/goleveldb/leveldb"

	"github.com/kellegous/golinks/pkg/store"
)

type Store struct {
	db *leveldb.DB
}

func Open(path string) (*Store, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Get(
	ctx context.Context,
	pattern *regexp.Regexp,
) (*store.Route, error) {
	key := keyFromPattern(pattern)
	val, err := s.db.Get(key, nil)
	if errors.Is(leveldb.ErrNotFound, err) {
		return nil, store.ErrorRouteNotfound
	} else if err != nil {
		return nil, err
	}

	return routeFromKeyAndVal(key, val)
}

func (s *Store) GetForURI(
	ctx context.Context,
	uri string,
) ([]*store.Route, error) {
	return nil, nil
}

func (s *Store) Put(
	ctx context.Context,
	r *store.Route,
) error {
	return s.db.Put(keyFromPattern(r.Pattern), valFromRoute(r), nil)
}
