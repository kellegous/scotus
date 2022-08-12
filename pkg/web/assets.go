package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist
var assets embed.FS

func getAssetsFS(dir string) (http.FileSystem, error) {
	if dir != "" {
		return http.Dir(dir), nil
	}

	s, err := fs.Sub(assets, "dist")
	if err != nil {
		return nil, err
	}

	return http.FS(s), nil
}
