package leveldb

import (
	"context"
	"errors"
	"regexp"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"

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
		return nil, store.ErrRouteNotfound
	} else if err != nil {
		return nil, err
	}

	return routeFromKeyAndVal(key, val)
}

func (s *Store) GetForPrefix(
	ctx context.Context,
	uri string,
) (store.Iterator[*store.Route], error) {
	key := keyFromString(uri)
	return routeIterator(
		s.db.NewIterator(
			&util.Range{Start: key, Limit: append(key, 0xff)},
			nil),
	), nil
}

func (s *Store) Put(
	ctx context.Context,
	r *store.Route,
) error {
	return s.db.Put(keyFromPattern(r.Pattern), valFromRoute(r), nil)
}

func rangeFromCursor(c string) (*util.Range, error) {
	if c == "" {
		return &util.Range{Start: nil, Limit: nil}, nil
	}

	cursor, err := decodeCursor(c)
	if err != nil {
		return nil, err
	}

	return &util.Range{Start: cursor, Limit: nil}, nil
}

func (s *Store) List(
	ctx context.Context,
	opts store.ListOptions,
) (store.Iterator[store.Cursor[*store.Route]], error) {
	r, err := rangeFromCursor(opts.Cursor)
	if err != nil {
		return nil, err
	}
	return cursorIterator(s.db.NewIterator(r, nil)), nil
}
