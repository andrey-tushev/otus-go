package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := NewLoggingResponseWriter(w)

		start := time.Now()
		next.ServeHTTP(lrw, r)
		latency := time.Now().Sub(start)

		message := fmt.Sprintf("http %s [%s] %s %s %s %d %d %s",
			r.RemoteAddr,
			start,
			r.Method, r.RequestURI, r.Proto,
			lrw.StatusCode,
			latency.Microseconds(),
			r.Header.Get("User-Agent"),
		)
		s.logger.Info(message)
	})
}

// Кастомный ResponseWriter, который позволяет вытащить http код.
// ASK: Можно ли решить это как то менее костыльно?
type loggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
