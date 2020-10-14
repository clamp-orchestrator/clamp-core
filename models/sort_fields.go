package models

import (
	"errors"
	"fmt"
	"strings"
)

//SortByFields contains an array of []sortBy type
//Used to store all sortBy key-value pairs
type SortByFields []sortBy

type sortBy struct {
	key   string
	order string
}

//ParseFromQuery is used to parse sortBy query string to a custom SortByFields type
//It returns an ordered sortBy struct containing KEY VALUE pair
//If an unknown key is used, an error is raised
func (sortArr *SortByFields) ParseFromQuery(sortByString string) error {
	if len(sortByString) == 0 {
		return nil
	}
	sortByString = cleanUpQuery(sortByString)
	sortByArgs := strings.Split(sortByString, ";")
	for _, value := range sortByArgs {
		sortPair := strings.Split(value, ":")

		if len(sortPair) != 2 || !verifySortValues(sortPair[0], sortPair[1]) {
			return errors.New("Unsupported value provided for sortBy query")
		}
		key := sortPair[0]
		value := sortPair[1]
		sort := sortBy{key: key, order: value}
		*sortArr = append(*sortArr, sort)
		fmt.Println(sortArr)
	}
	return nil
}

func cleanUpQuery(sortByQuery string) string {
	length := len(sortByQuery)
	sortByQuery = strings.ToLower(sortByQuery)
	if sortByQuery[length-1] == ';' {
		sortByQuery = sortByQuery[0 : length-1]
	}
	return sortByQuery
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
