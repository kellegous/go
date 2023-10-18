package web

import (
	"embed"
	"io/fs"
	"net/http"
	"net/http/httputil"
)

//go:embed ui
var ui embed.FS

func assetsHandler(opts *Options) (http.Handler, error) {
	f, err := fs.Sub(ui, "ui")
	if err != nil {
		return nil, err
	}

	if url := opts.assetProxyBaseURL; url != nil {
		return httputil.NewSingleHostReverseProxy(url), nil
	}

	return http.FileServer(http.FS(f)), nil
}
