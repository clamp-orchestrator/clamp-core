package parsers

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
)

//reference human readable keys to DB key values
var keyReferences = map[string]string{"id": "id", "createddate": "created_at", "name": "name"}

//SortByQueryParser is used to parse sortBy query JSON string to a map of (string,string)
//If any required key is missing, it is set as default
//All unknown keys are ignored
func SortByQueryParser(sortByQuery string) (map[string]string, error) {
	if sortByQuery == "" {
		sortByQuery = "{}"
	}
	sortByQuery = strings.ToLower(sortByQuery)
	sortBy := map[string]string{}
	err := json.Unmarshal([]byte(sortByQuery), &sortBy)
	var cleanedSortByArgs map[string]string
	if err == nil {
		cleanedSortByArgs, err = cleanSortByQuery(sortBy)
	}

	if err != nil {
		log.Println(err)
		return make(map[string]string), errors.New("Unsupported format for sortBy Query")
	}
	return cleanedSortByArgs, nil
}

func cleanSortByQuery(sortBy map[string]string) (map[string]string, error) {
	supportedKeys := []string{"id", "createddate", "name"}
	cleanedSortByArgs := make(map[string]string)
	for _, key := range supportedKeys {
		value, found := sortBy[key]
		if !found {
			continue
		}
		if value != "asc" && value != "desc" {
			return make(map[string]string), errors.New("Non supported argument for key " + key)
		}
		referenceKey, found := keyReferences[key]
		if found {
			cleanedSortByArgs[referenceKey] = value
		}
	}
	return cleanedSortByArgs, nil
}
