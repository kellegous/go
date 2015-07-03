package context

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

const (
	routesDbFilename        = "routes.db"
	idLogFilename           = "id"
	idBatchSize      uint64 = 1000
)

// Route ...
type Route struct {
	URL  string
	Time time.Time
}

//
func (o *Route) write(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, o.Time.UnixNano()); err != nil {
		return err
	}

	if _, err := w.Write([]byte(o.URL)); err != nil {
		return err
	}

	return nil
}

//
func (o *Route) read(r io.Reader) error {
	var t int64
	if err := binary.Read(r, binary.LittleEndian, &t); err != nil {
		return err
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	o.URL = string(b)
	o.Time = time.Unix(0, t)
	return nil
}

// Context ...
type Context struct {
	path string
	db   *leveldb.DB
	lck  sync.Mutex
	id   uint64
}

// Open ...
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

	c := &Context{
		path: path,
		db:   db,
	}

	// make sure we have an id log file
	if _, err := os.Stat(filepath.Join(path, idLogFilename)); err != nil {
		if err := c.commit(idBatchSize); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Get ...
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

// Put ...
func (c *Context) Put(key string, rt *Route) error {
	var buf bytes.Buffer
	if err := rt.write(&buf); err != nil {
		return err
	}

	return c.db.Put([]byte(key), buf.Bytes(), &opt.WriteOptions{Sync: true})
}

func (c *Context) commit(id uint64) error {
	w, err := os.Create(filepath.Join(c.path, idLogFilename))
	if err != nil {
		return err
	}
	defer w.Close()

	if err := binary.Write(w, binary.LittleEndian, id); err != nil {
		return err
	}

	return w.Sync()
}

// NextID ...
func (c *Context) NextID() (uint64, error) {
	c.lck.Lock()
	defer c.lck.Unlock()

	// when we hit a batch boundary, we will commit all ids until the next
	// boundary. If we crash, we'll just throw away a batch of ids in the worst
	// case.
	if c.id%idBatchSize == 0 {
		if err := c.commit(c.id + idBatchSize); err != nil {
			return 0, err
		}
	}

	c.id++

	return c.id, nil
}
