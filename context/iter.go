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

// Valid indicates whether the current values of the iterator are valid.
func (i *Iter) Valid() bool {
	return i.it.Valid() && i.err == nil
}

// Next advances the iterator to the next value.
func (i *Iter) Next() bool {
	it := i.it

	i.name = ""
	i.rt = nil

	if !it.Next() {
		return false
	}

	rt := &Route{}

	if err := rt.read(bytes.NewBuffer(it.Value())); err != nil {
		i.err = err
		return false
	}

	i.name = string(i.it.Key())
	i.rt = rt

	return true
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
