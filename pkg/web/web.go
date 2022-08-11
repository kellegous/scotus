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
	data *Data,
) error {
	m := http.NewServeMux()

	m.HandleFunc(
		"/api/debug/build",
		func(w http.ResponseWriter, r *http.Request) {
			ctx, done, _ := logging.ForRequest(ctx, time.Minute)
			defer done()
			sendJSONOK(ctx, w, data.Build)
		})

	return http.ListenAndServe(addr, m)
}
