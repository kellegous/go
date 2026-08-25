package leveldb

import (
	"context"
	"errors"
	"iter"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"google.golang.org/protobuf/proto"

	"github.com/kellegous/poop"

	"github.com/kellegous/golinks"
	"github.com/kellegous/golinks/internal/config"
	"github.com/kellegous/golinks/internal/store"
)

type Store struct {
	db *leveldb.DB
}

func Open(path string) (store.Store, error) {
	if path == "" {
		return nil, poop.New("path is required")
	}

	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func FromOptions(options iter.Seq2[*config.Option, error]) (store.Store, error) {
	var path string
	for option, err := range options {
		if err != nil {
			return nil, poop.Chain(err)
		}
		switch option.Key {
		case "path":
			path = option.Val
		}
	}
	return Open(path)
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Get(ctx context.Context, prefix string) (*golinks.Link, error) {
	return getLink(ctx, s.db, prefix)
}

type dbOrTx interface {
	Get(key []byte, opts *opt.ReadOptions) ([]byte, error)
	Put(key []byte, value []byte, opts *opt.WriteOptions) error
	Delete(key []byte, opts *opt.WriteOptions) error
	Has(key []byte, opts *opt.ReadOptions) (bool, error)
}

func getLink(ctx context.Context, db dbOrTx, prefix string) (*golinks.Link, error) {
	val, err := db.Get([]byte(prefix), nil)
	if errors.Is(err, leveldb.ErrNotFound) {
		return nil, poop.Chain(store.ErrLinkNotfound)
	} else if err != nil {
		return nil, poop.Chain(err)
	}

	var link golinks.Link
	if err := proto.Unmarshal(val, &link); err != nil {
		return nil, poop.Chain(err)
	}

	return &link, nil
}

func (s *Store) Put(ctx context.Context, link *golinks.Link) (*golinks.Link, error) {
	// TODO(kellegous): where should validation be done?
	prefix := link.GetPrefix()
	if prefix == "" {
		return nil, poop.New("prefix is required")
	}

	val, err := proto.Marshal(link)
	if err != nil {
		return nil, poop.Chain(err)
	}

	if err := s.db.Put([]byte(prefix), val, nil); err != nil {
		return nil, poop.Chain(err)
	}

	return link, nil
}

func (s *Store) Delete(ctx context.Context, prefix string) (*golinks.Link, error) {
	tx, err := s.db.OpenTransaction()
	if err != nil {
		return nil, poop.Chain(err)
	}
	defer tx.Discard()

	link, err := getLink(ctx, tx, prefix)
	if err != nil {
		return nil, poop.Chain(err)
	}

	if err := tx.Delete([]byte(prefix), nil); err != nil {
		return nil, poop.Chain(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, poop.Chain(err)
	}

	return link, nil
}

func (s *Store) List(ctx context.Context, start string) iter.Seq2[*golinks.Link, error] {
	return func(yield func(*golinks.Link, error) bool) {
		it := s.db.NewIterator(&util.Range{Start: []byte(start)}, nil)
		defer it.Release()

		for it.Next() {
			var link golinks.Link
			if err := proto.Unmarshal(it.Value(), &link); err != nil {
				yield(nil, poop.Chain(err))
				return
			}

			if !yield(&link, nil) {
				return
			}
		}

		if err := it.Error(); err != nil {
			yield(nil, err)
		}
	}
}
