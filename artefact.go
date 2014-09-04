package main

import (
	"encoding/json"
)

type Detail struct {
	NeedIDs             []string `json:"need_ids"`
	BusinessProposition bool     `json:"business_proposition"`
	Description         string   `json:"description"`
}

type Artefact struct {
	ID      string `json:"id"`
	WebURL  string `json:"web_url"`
	Title   string `json:"title"`
	Format  string `json:"format"`
	Details Detail `json:"details"`
}

func ParseArtefactResponse(response []byte) (*Artefact, error) {
	artefact := &Artefact{}
	if err := json.Unmarshal(response, &artefact); err != nil {
		return nil, err
	}

	return artefact, nil
}
