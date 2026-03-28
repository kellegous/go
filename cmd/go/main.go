package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/kellegous/glue/devmode"
	"github.com/kellegous/go/internal/backend"
	"github.com/kellegous/go/internal/backend/firestore"
	"github.com/kellegous/go/internal/backend/leveldb"
	"github.com/kellegous/go/internal/ui"
	"github.com/kellegous/go/internal/web"
)

func getAssets(ctx context.Context, devMode *devmode.Flag) (http.Handler, error) {
	if !devMode.IsEnabled() {
		return ui.Assets()
	}

	return devmode.AssetsFromVite(ctx, devMode)
}

func getBackend() (backend.Backend, error) {
	switch viper.GetString("backend") {
	case "leveldb":
		return leveldb.New(viper.GetString("data"))
	case "firestore":
		return firestore.New(context.Background(), viper.GetString("project"))
	default:
		return nil, fmt.Errorf("unknown backend %s", viper.GetString("backend"))
	}
}

func main() {
	var devMode devmode.Flag
	pflag.String("addr", ":8067", "default bind address")
	pflag.Bool("admin", false, "allow admin-level requests")
	pflag.String("version", "", "version string")
	pflag.String("backend", "leveldb", "backing store to use. 'leveldb' and 'firestore' currently supported.")
	pflag.String("data", "data", "The location of the leveldb data directory")
	pflag.String("project", "", "The GCP project to use for the firestore backend. Will attempt to use application default creds if not defined.")
	pflag.String("host", "", "The host field to use when gnerating the source URL of a link. Defaults to the Host header of the generate request")
	pflag.Var(
		&devMode,
		"dev-mode",
		"Enable dev mode")
	pflag.Parse()

	ctx := context.Background()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Panic(err)
	}

	// allow env vars to set pflags
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	backend, err := getBackend()
	if err != nil {
		log.Panic(err)
	}
	defer backend.Close()

	assets, err := getAssets(ctx, &devMode)
	if err != nil {
		log.Panic(err)
	}

	go func() {
		if err := devMode.ShowBannerWhenReady(ctx, os.Stdout, viper.GetString("addr")); err != nil {
			log.Panic(err)
		}
	}()

	log.Panic(web.ListenAndServe(backend, assets))
}
