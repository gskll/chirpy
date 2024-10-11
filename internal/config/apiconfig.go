package config

import (
	"sync/atomic"
)

type ApiConfig struct {
	FileServerHits atomic.Int32
}

func NewApiConfig() *ApiConfig {
	return &ApiConfig{}
}
