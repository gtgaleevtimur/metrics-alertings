package handler

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"time"
)

var (
	Log *zerolog.Logger
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.DurationFieldUnit = time.Millisecond
	consoleLog := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: time.DateTime + ".000",
	}
	l := zerolog.New(consoleLog).With().Timestamp().Caller().Logger()
	Log = &l
}

func SetLevel(s string) {
	v, err := zerolog.ParseLevel(s)
	if err != nil {
		Log.Warn().Err(err).Msg("invalid log level specified")
	}
	zerolog.SetGlobalLevel(v)
}

func LoggerHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// to get response data
		ww := middleware.NewWrapResponseWriter(writer, request.ProtoMajor)
		tStart := time.Now()
		defer func() {
			dur := time.Since(tStart)
			Log.Info().
				Str("host", request.Host).
				Str("url", request.URL.Path).
				Str("method", request.Method).
				Dur("ms_served", dur).
				Int("status", ww.Status()).
				Int("bytes_sent", ww.BytesWritten()).
				Msg("http request")
		}()
		next.ServeHTTP(ww, request)
	})
}
