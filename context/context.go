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
	routesDbFilename = "routes.db"
	idLogFilename    = "id"
)

// Route ...
type Route struct {
	URL  string    `json:"url"`
	Time time.Time `json:"time"`
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

	id, err := load(filepath.Join(path, idLogFilename))
	if err != nil {
		return nil, err
	}

	return &Context{
		path: path,
		db:   db,
		id:   id,
	}, nil
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

// Del ...
func (c *Context) Del(key string) error {
	return c.db.Delete([]byte(key), &opt.WriteOptions{Sync: true})
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

	c.id++

	if err := commit(filepath.Join(c.path, idLogFilename), c.id); err != nil {
		return 0, err
	}

	return c.id, nil
}
