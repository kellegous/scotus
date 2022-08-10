package web

import (
	"context"
	"net/http"
	"time"

	"github.com/kellegous/scotus/pkg/logging"
)

func ListenAndServe(
	ctx context.Context,
	addr string,
) error {
	m := http.NewServeMux()

	m.HandleFunc(
		"/api/debug/buildinfo",
		func(w http.ResponseWriter, r *http.Request) {
			ctx, done, _ := logging.ForRequest(ctx, time.Minute)
			defer done()
			apiGetBuildInfo(ctx, w, r)
		})

	return http.ListenAndServe(addr, m)
}
