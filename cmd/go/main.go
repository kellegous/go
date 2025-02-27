package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/kellegous/go/backend"
	"github.com/kellegous/go/backend/firestore"
	"github.com/kellegous/go/backend/leveldb"
	"github.com/kellegous/go/internal/ui"
	"github.com/kellegous/go/web"
)

func getAssets(proxyURL *url.URL) (http.Handler, error) {
	if proxyURL == nil {
		return ui.Assets()
	}
	p := httputil.NewSingleHostReverseProxy(proxyURL)
	dir := p.Director
	p.Director = func(r *http.Request) {
		dir(r)
		r.Host = proxyURL.Host
	}
	return p, nil
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

type URL struct {
	*url.URL
}

func (u *URL) Set(v string) error {
	var err error
	u.URL, err = url.Parse(v)
	return err
}

func (u *URL) Type() string {
	return "url"
}

func (u *URL) String() string {
	if u.URL == nil {
		return ""
	}
	return u.URL.String()
}

func main() {
	var assetProxyURL URL
	pflag.String("addr", ":8067", "default bind address")
	pflag.Bool("admin", false, "allow admin-level requests")
	pflag.String("version", "", "version string")
	pflag.String("backend", "leveldb", "backing store to use. 'leveldb' and 'firestore' currently supported.")
	pflag.String("data", "data", "The location of the leveldb data directory")
	pflag.String("project", "", "The GCP project to use for the firestore backend. Will attempt to use application default creds if not defined.")
	pflag.String("host", "", "The host field to use when gnerating the source URL of a link. Defaults to the Host header of the generate request")
	pflag.Var(
		&assetProxyURL,
		"asset-proxy-url",
		"The URL to use for proxying asset requests")
	pflag.Parse()

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

	assets, err := getAssets(assetProxyURL.URL)
	if err != nil {
		log.Panic(err)
	}

	log.Panic(web.ListenAndServe(backend, assets))
}
