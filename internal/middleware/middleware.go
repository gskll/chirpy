package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gskll/chirpy2/internal/config"
)

type Middleware struct {
	cfg *config.ApiConfig
}

func NewMiddleware(cfg *config.ApiConfig) *Middleware {
	return &Middleware{cfg: cfg}
}

func (m *Middleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func (m *Middleware) Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.cfg.FileServerHits.Add(1)

		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")

		next.ServeHTTP(w, r)
	})
}
