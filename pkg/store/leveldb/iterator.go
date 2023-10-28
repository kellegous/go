package leveldb

import (
	"context"

	"github.com/syndtr/goleveldb/leveldb/iterator"

	"github.com/kellegous/golinks/pkg/store"
)

type Iterator[T any] struct {
	it    iterator.Iterator
	fn    func(key []byte, val []byte) (T, error)
	empty T
}

func (i *Iterator[T]) Next(ctx context.Context) (T, error) {
	if !i.it.Next() {
		if err := i.it.Error(); err != nil {
			return i.empty, err
		}

		return i.empty, store.ErrIteratorDone
	}

	return i.fn(i.it.Key(), i.it.Value())
}

func (i *Iterator[T]) Close() error {
	i.it.Release()
	return nil
}

func routeIterator(it iterator.Iterator) *Iterator[*store.Route] {
	return &Iterator[*store.Route]{
		it: it,
		fn: func(key []byte, val []byte) (*store.Route, error) {
			return routeFromKeyAndVal(key, val)
		},
	}
}

func cursorIterator(it iterator.Iterator) *Iterator[store.Cursor[*store.Route]] {
	return &Iterator[store.Cursor[*store.Route]]{
		it: it,
		fn: func(key []byte, val []byte) (store.Cursor[*store.Route], error) {
			r, err := routeFromKeyAndVal(key, val)
			if err != nil {
				return nil, err
			}

			return &Cursor{r: r}, nil
		},
	}
}
