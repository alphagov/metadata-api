package content

import ()

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
