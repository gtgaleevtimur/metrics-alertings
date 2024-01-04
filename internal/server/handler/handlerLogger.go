package handler

import (
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var (
	sugar *zap.SugaredLogger
)

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	sugar = logger.Sugar()
}

func LoggerHandlerMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		logFn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			tStart := time.Now()
			defer func() {
				dur := time.Since(tStart)
				sugar.Infoln(
					"host", r.Host,
					"url", r.URL.Path,
					"method", r.Method,
					"status", ww.Status(), // получаем перехваченный код статуса ответа
					"duration", dur,
					"size", ww.BytesWritten(), // получаем перехваченный размер ответа
				)
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(logFn)
	}
}
