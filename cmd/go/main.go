package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/ctSkennerton/shortlinks/backend"
	"github.com/ctSkennerton/shortlinks/backend/firestore"
	"github.com/ctSkennerton/shortlinks/backend/leveldb"
	"github.com/ctSkennerton/shortlinks/web"
)

func main() {
	pflag.String("addr", ":8067", "default bind address")
	pflag.Bool("admin", false, "allow admin-level requests")
	pflag.String("version", "", "version string")
	pflag.String("backend", "leveldb", "backing store to use. 'leveldb' and 'firestore' currently supported.")
	pflag.String("data", "data", "The location of the leveldb data directory")
	pflag.String("project", "", "The GCP project to use for the firestore backend. Will attempt to use application default creds if not defined.")
	pflag.String("host", "", "The host field to use when gnerating the source URL of a link. Defaults to the Host header of the generate request")
	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Panic(err)
	}

	// allow env vars to set pflags
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	var backend backend.Backend

	switch viper.GetString("backend") {
	case "leveldb":
		var err error
		backend, err = leveldb.New(viper.GetString("data"))
		if err != nil {
			log.Panic(err)
		}
	case "firestore":
		var err error

		backend, err = firestore.New(context.Background(), viper.GetString("project"))
		if err != nil {
			log.Panic(err)
		}
	default:
		log.Panic(fmt.Sprintf("unknown backend %s", viper.GetString("backend")))
	}

	defer backend.Close()

	log.Panic(web.ListenAndServe(backend))
}
