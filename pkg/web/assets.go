package web

import (
	"embed"
	"io/fs"
	"net/http"
	"path/filepath"
)

//go:embed ui
var ui embed.FS

type Assets struct {
	paths map[string]bool
	http.Handler
}

func (a *Assets) CanServe(path string) bool {
	return a.paths[path]
}

func assetsIn(fsys fs.FS, dir string) (*Assets, error) {
	f, err := fs.Sub(fsys, dir)
	if err != nil {
		return nil, err
	}

	paths := map[string]bool{}
	if err := fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		p := filepath.Clean("/" + path)
		if d.Name() == "index.html" {
			d := filepath.Dir(p)
			paths[d] = true
			if d != "/" {
				paths[d+"/"] = true
			}
		} else {
			paths[filepath.Clean(p)] = true
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &Assets{
		paths:   paths,
		Handler: http.FileServer(http.FS(f)),
	}, nil
}
