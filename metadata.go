package main

import (
	"github.com/alphagov/metadata-api/content_api"
	"github.com/alphagov/metadata-api/need_api"
)

type ResponseInfo struct {
	Status string `json:"status"`
}

type Metadata struct {
	Artefact     *content_api.Artefact `json:"artefact"`
	Need         *need_api.Need        `json:"need"`
	ResponseInfo *ResponseInfo         `json:"_response_info"`
}
