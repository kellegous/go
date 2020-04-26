package leveldb

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/stgarf/go-links/internal"
)

const (
	routesDbFilename = "routes.db"
	idLogFilename    = "id"
)

// Backend provides access to the leveldb store.
type Backend struct {
	// Path contains the location on disk where this DB exists.
	path string
	db   *leveldb.DB
	lck  sync.Mutex
	id   uint64
}

// Commit the given ID to the data store.
func commit(filename string, id uint64) error {
	w, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer w.Close()

	if err := binary.Write(w, binary.LittleEndian, id); err != nil {
		return err
	}

	return w.Sync()
}

// Load the current ID from the data store.
func load(filename string) (uint64, error) {
	if _, err := os.Stat(filename); err != nil {
		return 0, commit(filename, 0)
	}

	r, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer r.Close()

	var id uint64
	if err := binary.Read(r, binary.LittleEndian, &id); err != nil {
		return 0, err
	}

	return id, nil
}

// New instantiates a new Backend
func New(path string) (*Backend, error) {
	backend := Backend{
		path: path,
	}

	if _, err := os.Stat(backend.path); err != nil {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return nil, err
		}
	}

	// open the database
	db, err := leveldb.OpenFile(filepath.Join(backend.path, routesDbFilename), nil)
	if err != nil {
		return nil, err
	}
	backend.db = db

	id, err := load(filepath.Join(backend.path, idLogFilename))
	if err != nil {
		return nil, err
	}
	backend.id = id

	return &backend, nil
}

// Close the resources associated with this backend.
func (backend *Backend) Close() error {
	return backend.db.Close()
}

// Get retreives a shortcut from the data store.
func (backend *Backend) Get(ctx context.Context, name string) (*internal.Route, error) {
	val, err := backend.db.Get([]byte(name), nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return nil, internal.ErrRouteNotFound
		}
		return nil, err
	}

	rt := &internal.Route{}
	if err := rt.Read(bytes.NewBuffer(val)); err != nil {
		return nil, err
	}

	return rt, nil
}

// Put stores a new shortcut in the data store.
func (backend *Backend) Put(ctx context.Context, key string, rt *internal.Route) error {
	var buf bytes.Buffer
	if err := rt.Write(&buf); err != nil {
		return err
	}

	return backend.db.Put([]byte(key), buf.Bytes(), &opt.WriteOptions{Sync: true})
}

// Del removes an existing shortcut from the data store.
func (backend *Backend) Del(ctx context.Context, key string) error {
	return backend.db.Delete([]byte(key), &opt.WriteOptions{Sync: true})
}

// List all routes in an iterator, starting with the key prefix of start (which can also be nil).
func (backend *Backend) List(ctx context.Context, start string) (internal.RouteIterator, error) {
	return &RouteIterator{
		it: backend.db.NewIterator(&util.Range{
			Start: []byte(start),
			Limit: nil,
		}, nil),
	}, nil
}

// GetAll gets everything in the db to dump it out for backup purposes
func (backend *Backend) GetAll(ctx context.Context) (map[string]internal.Route, error) {
	golinks := map[string]internal.Route{}
	iter := backend.db.NewIterator(nil, nil)
	defer iter.Release()

	for iter.Next() {
		key := iter.Key()
		val := iter.Value()
		rt := &internal.Route{}
		if err := rt.Read(bytes.NewBuffer(val)); err != nil {
			return nil, err
		}
		golinks[string(key[:])] = *rt
	}

	if err := iter.Error(); err != nil {
		return nil, err
	}

	return golinks, nil
}

func (backend *Backend) commit(id uint64) error {
	w, err := os.Create(filepath.Join(backend.path, idLogFilename))
	if err != nil {
		return err
	}
	defer w.Close()

	if err := binary.Write(w, binary.LittleEndian, id); err != nil {
		return err
	}

	return w.Sync()
}

// NextID generates the next numeric ID to be used for an auto-named shortcut.
func (backend *Backend) NextID(ctx context.Context) (uint64, error) {
	backend.lck.Lock()
	defer backend.lck.Unlock()

	backend.id++

	if err := commit(filepath.Join(backend.path, idLogFilename), backend.id); err != nil {
		return 0, err
	}

	return backend.id, nil
}
