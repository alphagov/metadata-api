package main

import (
	"os"
)

type Config struct {
	BearerTokenNeedAPI    string
}

func InitConfig() *Config {
	return &Config{
		BearerTokenNeedAPI:    os.Getenv("NEED_API_BEARER_TOKEN"),
	}
}
