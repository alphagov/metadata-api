package request

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	NotFoundError error = errors.New("not found")
)

func NewRequest(url, bearerToken string) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+bearerToken)
	request.Header.Add("Accept", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusNotFound {
		return nil, NotFoundError
	}

	return response, err
}

func ReadResponseBody(response *http.Response) (string, error) {
	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	return strings.TrimSpace(string(body)), err
}
