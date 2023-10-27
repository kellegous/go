package leveldb

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/kellegous/golinks/pkg/store"
)

const (
	version     byte = 0
	routePrefix byte = 'r'
)

func keyFromPattern(p *regexp.Regexp) []byte {
	s := p.String()
	b := make([]byte, len(s)+1)
	b[0] = routePrefix
	copy(b[1:], s)
	return b
}

func keyFromString(s string) []byte {
	b := make([]byte, len(s)+1)
	b[0] = routePrefix
	copy(b[1:], s)
	return b
}

func valFromRoute(r *store.Route) []byte {
	b := make([]byte, len(r.URL)+9)

	// byte index 0 is the version
	b[0] = version

	// byte indexes 1-8 are the timestamp
	var ts int64
	if !r.Time.IsZero() {
		ts = r.Time.Unix()
	}

	binary.BigEndian.PutUint64(b[1:], uint64(ts))

	// byte indexes 9+ are the URL
	copy(b[9:], r.URL)

	return b
}

func routeFromKeyAndVal(
	key []byte,
	val []byte,
) (*store.Route, error) {
	if len(val) < 10 {
		return nil, errors.New("invalid value: too short")
	}

	if val[0] != version {
		return nil, errors.New("invalid value: invalid version")
	}

	if key[0] != routePrefix {
		return nil, errors.New("invalid key: invalid prefix")
	}

	pattern, err := regexp.Compile(string(key[1:]))
	if err != nil {
		return nil, fmt.Errorf("invalid key: %w", err)
	}

	var t time.Time
	if ts := binary.BigEndian.Uint64(val[1:]); ts != 0 {
		t = time.Unix(int64(ts), 0)
	}

	return &store.Route{
		Pattern: pattern,
		URL:     string(val[9:]),
		Time:    t,
	}, nil
}

func encodeCursor(p *regexp.Regexp) string {
	return base64.URLEncoding.EncodeToString(
		append(keyFromPattern(p), 0xff))
}

func decodeCursor(s string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(s)
}
