package config

import (
	"sync/atomic"

	"github.com/gskll/chirpy2/internal/database"
)

type ApiConfig struct {
	FileServerHits atomic.Int32
	DbQueries      *database.Queries
}

func NewApiConfig(dbQueries *database.Queries) *ApiConfig {
	return &ApiConfig{DbQueries: dbQueries}
}
