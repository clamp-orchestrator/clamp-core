package models

import (
	"errors"
	"fmt"
	"strings"
)

//SortByFields contains an array of []sortBy type
//Used to store all sortBy key-value pairs
type SortByFields []sortBy

//SortBy contains a key and an order to be used to sort a DB query by
//Order supported is asc/desc
type sortBy struct {
	Key   string
	Order string
}

//ParseFromQuery is used to parse sortBy query string to a custom SortByFields type
//It returns an ordered sortBy struct containing KEY VALUE pair
//If an unknown key is used, an error is raised
//Fields are seperated using a comma and key values are sepearted using a colon
func ParseFromQuery(sortByString string) (SortByFields, error) {
	var sortArr SortByFields = []sortBy{}
	if len(sortByString) == 0 {
		return sortArr, nil
	}
	sortByString = cleanUpQuery(sortByString)
	sortByArgs := strings.Split(sortByString, ",")
	for _, value := range sortByArgs {
		sortPair := strings.Split(value, ":")

		if len(sortPair) != 2 || !verifySortValues(sortPair[0], sortPair[1]) {
			return SortByFields{}, errors.New("Unsupported value provided for sortBy query")
		}
		key := sortPair[0]
		value := sortPair[1]
		sort := sortBy{Key: key, Order: value}
		sortArr = append(sortArr, sort)
		fmt.Println(sortArr)
	}
	return sortArr, nil
}

func cleanUpQuery(sortByQuery string) string {
	length := len(sortByQuery)
	sortByQuery = strings.ToLower(sortByQuery)
	if sortByQuery[length-1] == ',' {
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
