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
	JWTSecret      string
	PolkaKey       string
}

func NewApiConfig(db *database.Queries, platform, jwtSecret, polkaKey string) *ApiConfig {
	return &ApiConfig{Db: db, Platform: platform, JWTSecret: jwtSecret, PolkaKey: polkaKey}
}
