package content_api

import (
	"encoding/json"
	"errors"

	"github.com/alphagov/metadata-api/request"
)

var ErrUnparsableArtefact = errors.New("Got JSON that doesn't look like an artefact")

type Part struct {
	WebURL string `json:"web_url"`
	Title  string `json:"title"`
}

type Detail struct {
	NeedIDs             []string `json:"need_ids"`
	BusinessProposition bool     `json:"business_proposition"`
	Description         string   `json:"description"`
	Parts               []Part   `json:"parts"`
}

type Artefact struct {
	ID      string `json:"id"`
	WebURL  string `json:"web_url"`
	Title   string `json:"title"`
	Format  string `json:"format"`
	Details Detail `json:"details"`
}

func FetchArtefact(contentAPI, bearerToken, slug string) (*Artefact, error) {
	artefactResponse, err := request.NewRequest(contentAPI+slug+".json", bearerToken)
	if err != nil {
		return nil, err
	}

	if artefactResponse.StatusCode/100 != 2 {
		return nil, errors.New("Content API error " + artefactResponse.Status)
	}

	artefactBody, err := request.ReadResponseBody(artefactResponse)
	if err != nil {
		return nil, err
	}

	artefact, err := ParseArtefactResponse([]byte(artefactBody))
	if err != nil {
		return nil, err
	}

	return artefact, nil
}

func ParseArtefactResponse(response []byte) (*Artefact, error) {
	artefact := &Artefact{}
	if err := json.Unmarshal(response, &artefact); err != nil {
		return nil, err
	}

	if len(artefact.ID) == 0 {
		return nil, ErrUnparsableArtefact
	}

	return artefact, nil
}
