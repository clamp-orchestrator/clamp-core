package utils

import (
	"encoding/json"
	"errors"
	"log"
)

//reference human readable keys to DB key values
var keyReferences = map[string]string{"id": "id", "createdDate": "created_at", "name": "name"}

//ParseFilters is used to parse filters form JSON string to a map of (string,string)
//If any required key is missing, it is set as default
//All unknown keys are ignored
func ParseFilters(filterString string) (map[string]string, error) {
	if filterString == "" {
		filterString = "{}"
	}
	filters := map[string]string{}
	err := json.Unmarshal([]byte(filterString), &filters)
	var cleanedFilters map[string]string
	if err == nil {
		cleanedFilters, err = cleanFilters(filters)
	}
	if err != nil {
		log.Println(err)
		return make(map[string]string), errors.New("Illegal filter syntax")
	}
	return cleanedFilters, nil
}

func cleanFilters(filters map[string]string) (map[string]string, error) {
	supportedKeys := []string{"id", "createdDate", "name"}
	cleanedFilters := make(map[string]string)
	for _, key := range supportedKeys {
		value := filters[key]
		if value != "asc" && value != "desc" && value != "" {
			return cleanedFilters, errors.New("Non supported filter argument")
		}
		referenceKey, found := keyReferences[key]
		if found {
			cleanedFilters[referenceKey] = value
		}
	}
	return cleanedFilters, nil
}
