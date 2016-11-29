package main

import (
	"os"
)

type Config struct {
	BearerTokenContentAPI string
	BearerTokenNeedAPI    string
}

func InitConfig() *Config {
	return &Config{
		BearerTokenContentAPI: os.Getenv("CONTENT_API_BEARER_TOKEN"),
		BearerTokenNeedAPI:    os.Getenv("NEED_API_BEARER_TOKEN"),
	}
}
