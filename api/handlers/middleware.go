package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/parviz-yu/digital-wallet/pkg/logger"
)

type ctxKey int8

const (
	ctxKeyUserID ctxKey = iota
)

const (
	userIDHeader = "X-UserId"
	digestHeader = "X-Digest"
)

var (
	ErrNoUserIDHeader       = errors.New("X-UserId header required")
	ErrNoXDigestHeader      = errors.New("X-Digest header required")
	ErrInvalidXDigestHeader = errors.New("invalid X-Digest header value")
)

func AuthMiddlewareUserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get(userIDHeader)
		if userID == "" {
			Error(w, r, http.StatusUnauthorized, ErrNoUserIDHeader)
			return
		}

		ctx := context.WithValue(r.Context(), ctxKeyUserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewMWLogger(log logger.LoggerI) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := logger.With(
			log,
			logger.String("component", "middleware/logger"),
		)

		log.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := logger.With(
				log,
				logger.String("method", r.Method),
				logger.String("path", r.URL.Path),
				logger.String("remote_addr", r.RemoteAddr),
				logger.String("user_agent", r.UserAgent()),
				logger.String("request_id", middleware.GetReqID(r.Context())),
			)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				entry.Info("request completed",
					logger.Int("status", ww.Status()),
					logger.Int("bytes", ww.BytesWritten()),
					logger.String("duration", time.Since(t1).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
