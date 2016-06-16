package content_store

import (
	"encoding/json"
	"fmt"
	"strings"

	. "github.com/alphagov/metadata-api/content"
	"github.com/alphagov/plek/go"
)

func GetArtefact(slug string, api JSONRequest) (*Artefact, error) {
	jsonResponse, err := getJSON(slug, api)
	if err != nil {
		return nil, err
	}
	return parseJSON(jsonResponse)
}

func getJSON(slug string, api JSONRequest) (string, error) {
	contentStoreBase := fmt.Sprintf("%s/content/", plek.FindURL("content-store"))
	url := contentStoreBase + slug
	json, err := api.GetJSON(url, "")
	if err != nil {
		return "", err
	}
	return json, err
}

func parseJSON(response string) (*Artefact, error) {
	artefact := &Artefact{}
	var jsonMap map[string]interface{}
	if err := json.Unmarshal([]byte(response), &jsonMap); err != nil {
		return nil, err
	}

	documentType := jsonMap["document_type"]
	if documentType != nil {
		if strings.Contains(documentType.(string), "placeholder") {
			return nil, StatusError{404}
		}
	}

	artefact.ID = jsonMap["content_id"].(string)
	artefact.Title = jsonMap["title"].(string)
	artefact.Format = jsonMap["format"].(string)
	artefact.WebURL = webURL(jsonMap["base_path"].(string))
	artefact.Details = unmarshalDetails(jsonMap)
	artefact.Details.Parts = unmarshalParts(jsonMap, *artefact)

	return artefact, nil
}

func webURL(basePath string) string {
	webroot, _ := plek.WebsiteRoot()
	return fmt.Sprintf("%s%s", webroot, basePath)
}

func unmarshalDetails(jsonMap map[string]interface{}) Detail {
	detail := Detail{}
	needIds := jsonMap["need_ids"].([]interface{})
	stringNeedIds := make([]string, len(needIds))

	for i := range needIds {
		stringNeedIds[i] = needIds[i].(string)
	}

	detail.NeedIDs = stringNeedIds
	detail.Description = jsonMap["description"].(string)

	return detail
}

func unmarshalParts(jsonMap map[string]interface{}, artefact Artefact) []Part {
	jsonDetails := jsonMap["details"].(map[string]interface{})

	jsonParts, ok := jsonDetails["parts"].([]interface{})

	if ok {
		parts := []Part{}

		for i := range jsonParts {
			jsonPart := jsonParts[i].(map[string]interface{})
			part := Part{}
			part.WebURL = fmt.Sprintf("%s/%s", artefact.WebURL, jsonPart["slug"].(string))
			part.Title = jsonPart["title"].(string)
			parts = append(parts, part)
		}
		return parts
	}
	return []Part{}
}
