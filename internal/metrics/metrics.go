package metrics

import (
	"fmt"
	"net/http"

	"github.com/gskll/chirpy2/internal/config"
)

type Metrics struct {
	cfg *config.ApiConfig
}

func NewMetrics(cfg *config.ApiConfig) *Metrics {
	return &Metrics{cfg: cfg}
}

func (m *Metrics) Count(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.cfg.FileServerHits.Add(1)

		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")

		next.ServeHTTP(w, r)
	})
}

func (m *Metrics) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", m.cfg.FileServerHits.Load())))
}

func (m *Metrics) Reset(w http.ResponseWriter, r *http.Request) {
	m.cfg.FileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
