package performance_platform

import (
	"encoding/json"
	"time"
)

type Data struct {
	HumanID          string  `json:"humanId"`
	PagePath         string  `json:"pagePath"`
	SearchKeyword    string  `json:"searchKeyword"`
	SearchUniques    float32 `json:"searchUniques"`
	SearchUniquesSum float32 `json:"searchUniques:sum"`
	TimeSpan         string  `json:"timeSpan"`
	Type             string  `json:"dataType"`

	// Underscore fields mean something in backdrop?
	ID        string    `json:"_id"`
	Count     float32   `json:"_count"`
	Timestamp time.Time `json:"_timestamp"`
}

type Backdrop struct {
	Data    []Data "json:`data`"
	Warning string "json:`warning`"
}

func ParseBackdropResponse(response []byte) (*Backdrop, error) {
	backdropResponse := &Backdrop{}
	if err := json.Unmarshal(response, &backdropResponse); err != nil {
		return nil, err
	}

	return backdropResponse, nil
}
