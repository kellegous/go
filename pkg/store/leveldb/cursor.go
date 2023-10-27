package leveldb

import (
	"encoding/json"

	"github.com/kellegous/golinks/pkg/store"
)

type Cursor struct {
	r *store.Route
}

func (c *Cursor) Route() *store.Route {
	return c.r
}

func (c *Cursor) Cursor() string {
	return encodeCursor(c.r.Pattern)
}

func (c *Cursor) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.r)
}
