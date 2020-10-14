package parsers

import (
	"errors"
	"strings"
)

//SortByQueryParser is used to parse sortBy query JSON string to a map of (string,string)
//If any required key is missing, it is set as default
//All unknown keys are ignored
//Returns
func SortByQueryParser(sortByQuery string) (map[string]string, []string, error) {
	cleanedSortByArgs := make(map[string]string)
	sortOrder := []string{}
	length := len(sortByQuery)
	if length == 0 {
		return cleanedSortByArgs, sortOrder, nil
	}
	sortByQuery = strings.ToLower(sortByQuery)
	if sortByQuery[length-1] == ';' {
		sortByQuery = sortByQuery[0 : length-1]
	}
	sortByArgs := strings.Split(sortByQuery, ";")
	for _, value := range sortByArgs {
		sortPair := strings.Split(value, ":")

		if len(sortPair) != 2 || !verifySortValues(sortPair[0], sortPair[1]) {
			return map[string]string{}, []string{}, errors.New("Unsupported value provided for sortBy query")
		}
		key := sortPair[0]
		value := sortPair[1]
		cleanedSortByArgs[key] = value
		sortOrder = append(sortOrder, key)
	}
	return cleanedSortByArgs, sortOrder, nil
}

func verifySortValues(key string, value string) bool {
	if value != "desc" && value != "asc" {
		return false
	}
	supportedKeys := []string{"id", "name", "createdate"}
	for _, keyVal := range supportedKeys {
		if key == keyVal {
			return true
		}
	}
	return false
}
