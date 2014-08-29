package main

import (
	"net/http"
)

func NewRequest(url, bearerToken string) (*http.Response, error) {
	client := &http.Client{}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+bearerToken)
	return client.Do(request)
}
