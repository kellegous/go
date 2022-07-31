package firestore

import (
	"context"
	"time"

	fs "cloud.google.com/go/firestore"
	"github.com/ctSkennerton/shortlinks/internal"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NextID is the next numeric ID to use for auto-generated IDs
type NextID struct {
	ID uint32 `json:"id" firestore:"id"`
}

// Backend provides access to Google Firestore.
type Backend struct {
	db *fs.Client
}

// New instantiates a new Backend
func New(ctx context.Context, path string) (*Backend, error) {
	if path == "" {
		path = getGoogleProject()
	}
	client, err := fs.NewClient(ctx, path)
	if err != nil {
		return nil, err
	}
	backend := Backend{
		db: client,
	}

	return &backend, nil
}

// Close the resources associated with this backend.
func (backend *Backend) Close() error {
	return backend.db.Close()
}

// Get retreives a shortcut from the data store.
func (backend *Backend) Get(ctx context.Context, name string) (*internal.Route, error) {
	ref := backend.db.Doc("routes/" + name)

	snap, err := ref.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, internal.ErrRouteNotFound
		}
		return nil, err
	}

	var rt internal.Route
	if err := snap.DataTo(&rt); err != nil {
		return nil, err
	}

	return &rt, nil
}

// Put stores a new shortcut in the data store.
func (backend *Backend) Put(ctx context.Context, key string, rt *internal.Route) error {
	ref := backend.db.Doc("routes/" + key)

	_, err := ref.Set(ctx, rt)
	if err != nil {
		return err
	}

	return nil
}

// Del removes an existing shortcut from the data store.
func (backend *Backend) Del(ctx context.Context, key string) error {
	ref := backend.db.Doc("routes/" + key)

	_, err := ref.Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

// List all routes in an iterator, starting with the key prefix of start (which can also be nil).
func (backend *Backend) List(ctx context.Context, start string) (internal.RouteIterator, error) {
	col := backend.db.Collection("routes").OrderBy(fs.DocumentID, fs.Asc)

	if start != "" {
		// we have a starting ID.
		col = col.StartAt(start)
	}

	return &RouteIterator{
		ctx: ctx,
		db:  backend.db,
		it:  col.Documents(ctx),
	}, nil
}

// GetAll gets everything in the db to dump it out for backup purposes
func (backend *Backend) GetAll(ctx context.Context) (map[string]internal.Route, error) {
	golinks := map[string]internal.Route{}
	col := backend.db.Collection("routes").OrderBy(fs.DocumentID, fs.Asc)

	routes, err := col.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	for _, doc := range routes {
		var rt internal.Route
		if err := doc.DataTo(&rt); err != nil {
			return nil, err
		}
		golinks[doc.Ref.ID] = rt
	}
	return golinks, nil
}

// NextID generates the next numeric ID to be used for an auto-named shortcut.
func (backend *Backend) NextID(ctx context.Context) (uint64, error) {
	ref := backend.db.Doc("IDs/nextID")
	var nid uint32

	err := backend.db.RunTransaction(ctx, func(ctx context.Context, tx *fs.Transaction) error {
		var nextID *NextID

		doc, err := tx.Get(ref)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				// this is the very first auto-generated ID, we can make it
				// as :1 and return it
				nextID = new(NextID)
				nextID.ID = 1
				nid = 1
				err := tx.Create(ref, nextID)
				if err != nil {
					return err
				}
				return nil
			}
			return err
		}

		if err := doc.DataTo(&nextID); err != nil {
			return err
		}
		nextID.ID += 1
		nid = nextID.ID

		return tx.Set(ref, &nextID)
	})
	if err != nil {
		return 0, err
	}

	return uint64(nid), nil
}

func getGoogleProject() string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	creds, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return ""
	}
	return creds.ProjectID
}
