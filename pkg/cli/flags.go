package cli

import (
	"context"
	"fmt"

	"github.com/kellegous/golinks/pkg/backend"
	"github.com/kellegous/golinks/pkg/backend/firestore"
	"github.com/kellegous/golinks/pkg/backend/leveldb"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultBackend     = "leveldb"
	defaultHttpAddr    = ":8080"
	defaultLevelDBData = "data"
)

type firestoreFlags struct{}

func (f *firestoreFlags) Register(fs *pflag.FlagSet) {
	fs.String(
		"firestore.project",
		"",
		"The GCP project to use for the firestore backend. Will attempt to use application default creds if not defined.")
}

func (f *firestoreFlags) Project() string {
	return viper.GetString("firestore.project")
}

type leveldbFlags struct{}

func (f *leveldbFlags) Register(fs *pflag.FlagSet) {
	fs.String(
		"leveldb.data",
		defaultLevelDBData,
		"The location of the leveldb data directory")
}

func (f *leveldbFlags) Data() string {
	return viper.GetString("leveldb.data")
}

type withBackend struct {
	Firestore firestoreFlags
	LevelDB   leveldbFlags
}

func (f *withBackend) Backend(ctx context.Context) (backend.Backend, error) {
	switch b := viper.GetString("backend"); b {
	case "leveldb":
		return leveldb.New(f.LevelDB.Data())
	case "firestore":
		return firestore.New(ctx, f.Firestore.Project())
	default:
		return nil, fmt.Errorf("unknown backend %s", b)
	}
}

func (f *withBackend) Register(fs *pflag.FlagSet) {
	f.Firestore.Register(fs)
	f.LevelDB.Register(fs)
	// TODO(knorton): create supported string dynamically.
	fs.String(
		"backend",
		defaultBackend,
		"backing store to use. 'leveldb' and 'firestore' currently supported.")
}

type withHTTP struct{}

func (f *withHTTP) Register(fs *pflag.FlagSet) {
	fs.String(
		"http.host",
		"",
		"The host field to use when gnerating the source URL of a link. Defaults to the Host header of the generate request")
	fs.String(
		"http.addr",
		defaultHttpAddr,
		"The address to which the http server will bind")
}

func (f *withHTTP) Addr() string {
	return viper.GetString("http.addr")
}

func (f *withHTTP) Host() string {
	return viper.GetString("http.host")
}
