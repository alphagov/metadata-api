package content

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type JSONRequest interface {
	GetJSON(url string, bearerToken string) (string, error)
}

type StatusError struct {
	StatusCode int
}

func (e StatusError) Error() string {
	return fmt.Sprintf("The response had a %d status code", e.StatusCode)
}

type ApiRequest struct{}

func (api ApiRequest) GetJSON(url string, bearerToken string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	if resp.StatusCode/100 != 2 {
		return "", StatusError{resp.StatusCode}
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return string(body), err
}
