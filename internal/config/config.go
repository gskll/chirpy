package config

import (
	"sync/atomic"

	"github.com/gskll/chirpy2/internal/database"
)

const DEV = "dev"

type ApiConfig struct {
	FileServerHits atomic.Int32
	Db             *database.Queries
	Platform       string
}

func NewApiConfig(db *database.Queries, platform string) *ApiConfig {
	return &ApiConfig{Db: db, Platform: platform}
}
