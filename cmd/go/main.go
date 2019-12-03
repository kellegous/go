package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/kellegous/go/backend"
	"github.com/kellegous/go/backend/leveldb"
	"github.com/kellegous/go/web"
)

func main() {
	pflag.String("addr", ":8067", "default bind address")
	pflag.Bool("admin", false, "allow admin-level requests")
	pflag.String("version", "", "version string")
	pflag.String("backend", "leveldb", "backing store to use. Only 'leveldb' currently supported.")
	pflag.String("data", "data", "The location of the leveldb data directory")

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
	default:
		log.Panic(fmt.Sprintf("unknown backend %s", viper.GetString("backend")))
	}

	defer backend.Close()

	log.Panic(web.ListenAndServe(backend))
}
