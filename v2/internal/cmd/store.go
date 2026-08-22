package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/kellegous/poop"

	"github.com/kellegous/golinks/internal/config"
	"github.com/kellegous/golinks/internal/store"
	"github.com/kellegous/golinks/internal/store/leveldb"
)

type Store struct {
	StoreType StoreType
	Config    string
}

func (s *Store) Set(v string) error {
	t, opts, _ := strings.Cut(v, ":")
	st, err := getType(t)
	if err != nil {
		return poop.Chain(err)
	}
	s.StoreType = st
	s.Config = opts
	return nil
}

func (s *Store) String() string {
	return fmt.Sprintf("%s:%s", s.StoreType, s.Config)
}

func (s *Store) Type() string {
	return "type:options"
}

func (s *Store) Open(ctx context.Context) (store.Store, error) {
	switch s.StoreType {
	case StoreTypeLevelDB:
		return leveldb.FromOptions(
			config.WithDefaults(
				config.Parse(s.Config),
				config.Opt("path", "golinks.db"),
			),
		)
	}
	return nil, poop.Newf("unknown store type: %q", s.StoreType)
}

type StoreType string

const (
	StoreTypeLevelDB StoreType = "leveldb"
	StoreTypeSQLite  StoreType = "sqlite"
)

func getType(v string) (StoreType, error) {
	switch strings.ToLower(v) {
	case "leveldb":
		return StoreTypeLevelDB, nil
	case "sqlite":
		return StoreTypeSQLite, nil
	}
	return "", poop.Newf("unknown store type: %q", v)
}

func (s StoreType) String() string {
	return string(s)
}
