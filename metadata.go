package main

import (
	"github.com/alphagov/metadata-api/need_api"
	"github.com/alphagov/metadata-api/performance_platform"
)

type ResponseInfo struct {
	Status string `json:"status"`
}

type Metadata struct {
	Artefact     interface{}                      `json:"artefact"`
	Needs        []*need_api.Need                 `json:"needs"`
	Performance  *performance_platform.Statistics `json:"performance"`
	ResponseInfo *ResponseInfo                    `json:"_response_info"`
}
