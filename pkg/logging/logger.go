package logging

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"go.uber.org/zap"
)

type contextKey string

const loggerKey = contextKey("logger")

func L(ctx context.Context) *zap.Logger {
	if ctx != nil {
		if l, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
			return l
		}
	}
	return zap.L()
}

func ForRequest(
	ctx context.Context,
	timeout time.Duration,
) (context.Context, context.CancelFunc, *zap.Logger) {
	lg := L(ctx).With(zap.String("req_id", requestID()))
	ctx = context.WithValue(ctx, loggerKey, lg)
	ctx, done := context.WithTimeout(ctx, timeout)
	return ctx, done, lg
}

func requestID() string {
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(buf[:])
}
