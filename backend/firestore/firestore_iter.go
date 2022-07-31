package firestore

import (
	"context"
	"errors"

	fs "cloud.google.com/go/firestore"
	"github.com/ctSkennerton/shortlinks/internal"
	"google.golang.org/api/iterator"
)

// RouteIterator allows iteration of the named routes in firestore.
type RouteIterator struct {
	ctx context.Context
	db  *fs.Client
	it  *fs.DocumentIterator
	doc *fs.DocumentSnapshot
	err error
}

// Valid indicates whether the current values of the iterator are valid.
func (i *RouteIterator) Valid() bool {
	return i.db != nil && i.it != nil && i.doc != nil && i.err == nil
}

// Next advances the iterator to the next value.
func (i *RouteIterator) Next() bool {
	doc, err := i.it.Next()
	if err != nil {
		if errors.Is(err, iterator.Done) {
			i.doc = nil
		} else {
			i.err = err
		}
		return false
	}

	i.doc = doc

	return true
}

// Seek ...
func (i *RouteIterator) Seek(cur string) bool {
	// firestore makes this a little hard. Make a whole new
	// document iterator that starts at a new spot.
	i.it = i.db.Collection("routes").OrderBy(fs.DocumentID, fs.Asc).StartAt(cur).Documents(i.ctx)

	doc, err := i.it.Next()
	if err != nil {
		i.err = err
		i.doc = nil
		return false
	}

	i.doc = doc
	return true
}

// Error returns any active error that has stopped the iterator.
func (i *RouteIterator) Error() error {
	return i.err
}

// Name is the name of the current route.
func (i *RouteIterator) Name() string {
	return i.doc.Ref.ID
}

// Route is the current route.
func (i *RouteIterator) Route() *internal.Route {
	var rt internal.Route
	if err := i.doc.DataTo(&rt); err != nil {
		i.err = err
		return nil
	}
	return &rt
}

// Release disposes of the resources in the iterator.
func (i *RouteIterator) Release() {
	i.it.Stop()
	i.doc = nil
	i.err = nil
}
