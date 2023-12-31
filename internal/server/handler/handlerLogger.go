package handler

import (
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

var (
	sLog *zap.SugaredLogger
	once sync.Once
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func getLogger() *zap.SugaredLogger {
	once.Do(func() {
		logger, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
		sLog = logger.Sugar()
	})
	return sLog
}

func LoggerHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: writer, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}
		next.ServeHTTP(&lw, request) // внедряем реализацию http.ResponseWriter

		duration := time.Since(start)

		sugar := getLogger()

		sugar.Infoln(
			"uri", request.RequestURI,
			"method", request.Method,
			"status", responseData.status, // получаем перехваченный код статуса ответа
			"duration", duration,
			"size", responseData.size, // получаем перехваченный размер ответа
		)
	})
}
