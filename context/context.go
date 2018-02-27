package context

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const (
	routesDbFilename = "routes.db"
)

// Route is the value part of a shortcut.
type Route struct {
	URL       string    `json:"url"`
	Time      time.Time `json:"time"`
	Uid       uint64    `json:"uid"`
	Generated bool      `json:"generated"`
}

// Serialize this Route into the given Writer.
func (o *Route) write(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, o.Time.UnixNano()); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, o.Uid); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, o.Generated); err != nil {
		return err
	}

	if _, err := w.Write([]byte(o.URL)); err != nil {
		return err
	}

	return nil
}

// Deserialize this Route from the given Reader.
func (o *Route) read(r io.Reader) error {
	var t int64
	if err := binary.Read(r, binary.LittleEndian, &t); err != nil {
		return err
	}
	o.Time = time.Unix(0, t)

	if err := binary.Read(r, binary.LittleEndian, &o.Uid); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &o.Generated); err != nil {
		return err
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	o.URL = string(b)

	return nil
}

// Context provides access to the data store.
type Context struct {
	path string
	db   *leveldb.DB
}

// Open the context using path as the data store location.
func Open(path string) (*Context, error) {
	if _, err := os.Stat(path); err != nil {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return nil, err
		}
	}

	// open the database
	db, err := leveldb.OpenFile(filepath.Join(path, routesDbFilename), nil)
	if err != nil {
		return nil, err
	}

	return &Context{
		path: path,
		db:   db,
	}, nil
}

// Close the resources associated with this context.
func (c *Context) Close() error {
	return c.db.Close()
}

// Get retreives a shortcut from the data store.
func (c *Context) Get(name string) (*Route, error) {
	val, err := c.db.Get([]byte(name), nil)
	if err != nil {
		return nil, err
	}

	rt := &Route{}
	if err := rt.read(bytes.NewBuffer(val)); err != nil {
		return nil, err
	}

	return rt, nil
}

// Put stores a new shortcut in the data store.
func (c *Context) Put(key string, rt *Route) error {
	var buf bytes.Buffer
	if err := rt.write(&buf); err != nil {
		return err
	}

	return c.db.Put([]byte(key), buf.Bytes(), &opt.WriteOptions{Sync: true})
}

// Del removes an existing shortcut from the data store.
func (c *Context) Del(key string) error {
	return c.db.Delete([]byte(key), &opt.WriteOptions{Sync: true})
}

// List all routes in an iterator, starting with the key prefix of start (which can also be nil).
func (c *Context) List(start []byte) *Iter {
	return &Iter{
		it: c.db.NewIterator(&util.Range{
			Start: start,
			Limit: nil,
		}, nil),
	}
}

// GetAll gets everything in the db to dump it out for backup purposes
func (c *Context) GetAll() (map[string]Route, error) {
	golinks := map[string]Route{}
	iter := c.db.NewIterator(nil, nil)
	defer iter.Release()

	for iter.Next() {
		key := iter.Key()
		val := iter.Value()
		rt := &Route{}
		if err := rt.read(bytes.NewBuffer(val)); err != nil {
			return nil, err
		}
		golinks[string(key[:])] = *rt
	}

	if err := iter.Error(); err != nil {
		return nil, err
	}

	return golinks, nil
}
