package web

import (
	"context"
	"net/http"
	"time"
)

func ListenAndServe(
	ctx context.Context,
	addr string,
	assetsDir string,
	data *Data,
) error {
	m := NewMux(ctx)

	fs, err := getAssetsFS(assetsDir)
	if err != nil {
		return err
	}

	m.Handle("/", http.FileServer(fs))

	m.HandleFunc(
		"/api/debug/build",
		func(w http.ResponseWriter, r *http.Request) {
			ctx, done := context.WithTimeout(
				ContextFrom(w),
				time.Minute)
			defer done()
			sendJSONOK(ctx, w, data.Build)
		})

	return http.ListenAndServe(addr, m)
}
