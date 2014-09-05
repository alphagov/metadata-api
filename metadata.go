package main

import (
	"encoding/json"

	"github.com/alphagov/metadata-api/content_api"
	"github.com/alphagov/metadata-api/need_api"
)

type ResponseInfo struct {
	Status string `json:"status"`
}

type Metadata struct {
	Artefact     *content_api.Artefact `json:"artefact"`
	Needs        []*need_api.Need      `json:"need"`
	ResponseInfo *ResponseInfo         `json:"_response_info"`
}

func ParseMetadataResponse(response []byte) (*Metadata, error) {
	metadata := &Metadata{}
	if err := json.Unmarshal(response, &metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}
