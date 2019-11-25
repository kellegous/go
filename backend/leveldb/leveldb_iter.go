package leveldb

import (
	"bytes"

	"github.com/kellegous/go/internal"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

// Iter allows iteration of the named routes in the store.
type LevelDBRouteIterator struct {
	it   iterator.Iterator
	name string
	rt   *internal.Route
	err  error
}

func (i *LevelDBRouteIterator) decode() error {
	rt := &internal.Route{}
	if err := rt.Read(bytes.NewBuffer(i.it.Value())); err != nil {
		return err
	}

	i.name = string(i.it.Key())
	i.rt = rt
	return nil
}

// Valid indicates whether the current values of the iterator are valid.
func (i *LevelDBRouteIterator) Valid() bool {
	return i.it.Valid() && i.err == nil
}

// Next advances the iterator to the next value.
func (i *LevelDBRouteIterator) Next() bool {
	i.name = ""
	i.rt = nil

	if !i.it.Next() {
		return false
	}

	if err := i.decode(); err != nil {
		i.err = err
		return false
	}

	return true
}

// Seek ...
func (i *LevelDBRouteIterator) Seek(cur string) bool {
	i.name = ""
	i.rt = nil

	v := i.it.Seek([]byte(cur))

	if !i.it.Valid() {
		return v
	}

	if err := i.decode(); err != nil {
		i.err = err
	}

	return v
}

// Error returns any active error that has stopped the iterator.
func (i *LevelDBRouteIterator) Error() error {
	if err := i.it.Error(); err != nil {
		return err
	}

	return i.err
}

// Name is the name of the current route.
func (i *LevelDBRouteIterator) Name() string {
	return i.name
}

// Route is the current route.
func (i *LevelDBRouteIterator) Route() *internal.Route {
	return i.rt
}

// Release disposes of the resources in the iterator.
func (i *LevelDBRouteIterator) Release() {
	i.it.Release()
}
