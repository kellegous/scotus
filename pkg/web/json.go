package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kellegous/scotus/pkg/logging"

	"go.uber.org/zap"
)

func sendJSON(
	ctx context.Context,
	w http.ResponseWriter,
	status int,
	data interface{},
) {
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logging.L(ctx).Panic("unable to send JSON",
			zap.Error(err))
	}
}

func sendJSONOK(
	ctx context.Context,
	w http.ResponseWriter,
	data interface{},
) {
	sendJSON(ctx, w, http.StatusOK, data)
}

func sendJSONErr(
	ctx context.Context,
	w http.ResponseWriter,
	status int,
	msg string,
) {
	sendJSON(ctx, w, status, struct {
		Error string `json:"error"`
	}{
		Error: msg,
	})
}

func sendJSONServerErr(
	ctx context.Context,
	w http.ResponseWriter,
	err error,
) {
	sendJSONErr(ctx, w, http.StatusInternalServerError, "backend error")
	logging.L(ctx).Error("backend error",
		zap.Error(err))
}
