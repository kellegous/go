package context

import (
	"bytes"

	"github.com/syndtr/goleveldb/leveldb/iterator"
)

// Iter allows iteration of the named routes in the store.
type Iter struct {
	it   iterator.Iterator
	name string
	rt   *Route
	err  error
}

func (i *Iter) decode() error {
	rt := &Route{}
	if err := rt.read(bytes.NewBuffer(i.it.Value())); err != nil {
		return err
	}

	i.name = string(i.it.Key())
	i.rt = rt
	return nil
}

// Valid indicates whether the current values of the iterator are valid.
func (i *Iter) Valid() bool {
	return i.it.Valid() && i.err == nil
}

// Next advances the iterator to the next value.
func (i *Iter) Next() bool {
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
func (i *Iter) Seek(cur []byte) bool {
	i.name = ""
	i.rt = nil

	v := i.it.Seek(cur)

	if !i.it.Valid() {
		return v
	}

	if err := i.decode(); err != nil {
		i.err = err
	}

	return v
}

// Error returns any active error that has stopped the iterator.
func (i *Iter) Error() error {
	if err := i.it.Error(); err != nil {
		return err
	}

	return i.err
}

// Name is the name of the current route.
func (i *Iter) Name() string {
	return i.name
}

// Route is the current route.
func (i *Iter) Route() *Route {
	return i.rt
}

// Release disposes of the resources in the iterator.
func (i *Iter) Release() {
	i.it.Release()
}
