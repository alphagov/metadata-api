package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	BearerTokenContentAPI string `json:"bearer_token_content_api"`
	BearerTokenNeedAPI    string `json:"bearer_token_need_api"`
}

func ReadConfig(filename string) (*Config, error) {
	workingDirectory, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, err
	}

	configPath := workingDirectory + "/" + filename
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return config, nil
}
