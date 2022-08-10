package web

import (
	"context"
	"debug/buildinfo"
	"net/http"
	"os"
)

func apiGetBuildInfo(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) {
	bi, err := buildinfo.ReadFile(os.Args[0])
	if err != nil {
		sendJSONServerErr(ctx, w, err)
	}
	sendJSONOK(ctx, w, bi)
}
