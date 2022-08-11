package web

import (
	"context"
	"net/http"

	"github.com/kellegous/scotus/pkg/logging"
	"go.uber.org/zap"
)

type Mux struct {
	*http.ServeMux
	ctx context.Context
}

func NewMux(ctx context.Context) *Mux {
	return &Mux{
		ServeMux: http.NewServeMux(),
		ctx:      ctx,
	}
}

func (m *Mux) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx, lg := logging.ForRequest(m.ctx)
	rw := responseWriter{
		ResponseWriter: w,
		status:         http.StatusOK,
		ctx:            ctx,
	}

	m.ServeMux.ServeHTTP(&rw, r)
	lg.Info("http request",
		zap.Int("status", rw.status),
		zap.String("uri", r.RequestURI),
		zap.String("host", r.Host),
		zap.String("user-agent", r.UserAgent()))
}

type responseWriter struct {
	http.ResponseWriter
	ctx    context.Context
	status int
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func ContextFrom(w http.ResponseWriter) context.Context {
	if rw, ok := w.(*responseWriter); ok {
		return rw.ctx
	}
	return context.Background()
}
