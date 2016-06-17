package content_api

import (
	"encoding/json"

	. "github.com/alphagov/metadata-api/content"
	"github.com/alphagov/metadata-api/request"
)

func FetchArtefact(contentAPI, bearerToken, slug string) (*Artefact, error) {
	artefactResponse, err := request.NewRequest(contentAPI+slug+".json", bearerToken)
	if err != nil {
		return nil, err
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

	return artefact, nil
}
