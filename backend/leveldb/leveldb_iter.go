package leveldb

import (
	"bytes"

	"github.com/stgarf/go-links/internal"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

// RouteIterator allows iteration of the named routes in the store.
type RouteIterator struct {
	it   iterator.Iterator
	name string
	rt   *internal.Route
	err  error
}

func (i *RouteIterator) decode() error {
	rt := &internal.Route{}
	if err := rt.Read(bytes.NewBuffer(i.it.Value())); err != nil {
		return err
	}

	i.name = string(i.it.Key())
	i.rt = rt
	return nil
}

// Valid indicates whether the current values of the iterator are valid.
func (i *RouteIterator) Valid() bool {
	return i.it.Valid() && i.err == nil
}

// Next advances the iterator to the next value.
func (i *RouteIterator) Next() bool {
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
func (i *RouteIterator) Seek(cur string) bool {
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
func (i *RouteIterator) Error() error {
	if err := i.it.Error(); err != nil {
		return err
	}

	return i.err
}

// Name is the name of the current route.
func (i *RouteIterator) Name() string {
	return i.name
}

// Route is the current route.
func (i *RouteIterator) Route() *internal.Route {
	return i.rt
}

// Release disposes of the resources in the iterator.
func (i *RouteIterator) Release() {
	i.it.Release()
}
