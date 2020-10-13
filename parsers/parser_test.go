package parsers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortByParser(t *testing.T) {
	sortByString := `
	{
		"name": "desc",
		"id": "asc",
		"createddate": "desc"
	}
	`
	data, err := SortByQueryParser(sortByString)
	assert.Nil(t, err)
	assert.Equal(t, "asc", data["id"])
	assert.Equal(t, "desc", data["createddate"])
	assert.Equal(t, "desc", data["name"])
}

func TestSortByParserIgnoreUnknownFields(t *testing.T) {
	sortByString := `
	{
		"id": "asc",
		"createddate": "desc",
		"name": "desc",
		"not_not":"asc"
	}
	`
	data, err := SortByQueryParser(sortByString)
	assert.Nil(t, err)
	assert.Equal(t, "asc", data["id"])
	assert.Equal(t, "desc", data["createddate"])
	assert.Equal(t, "desc", data["name"])
	assert.NotContains(t, data, "not_not")
}

func TestSortByParserAssignBlankToMissing(t *testing.T) {
	sortByString := `
	{
		"createddate": "desc",
		"name": "desc"
	}
	`
	data, err := SortByQueryParser(sortByString)
	assert.Nil(t, err)
	assert.Equal(t, "", data["id"])
	assert.Equal(t, "desc", data["createddate"])
	assert.Equal(t, "desc", data["name"])
}

func TestSortByParserThrowErrorOnIllegalValue(t *testing.T) {
	sortByString := `
	{	
		"id":"randomValue",
		"createddate": "desc",
		"name": "desc"
	}
	`
	_, err := SortByQueryParser(sortByString)
	fmt.Println(err)
	assert.NotNil(t, err)
}

func TestSortByParserThrowErrorOnWrongJSON(t *testing.T) {
	sortByString := `
	{	
		"id","randomValue",
		"createddate": "desc",
		"name": "desc"
	}
	`
	_, err := SortByQueryParser(sortByString)
	fmt.Println(err)
	assert.NotNil(t, err)
}
func TestSortByParserAllowEmptyJSON(t *testing.T) {
	sortByString := "{}"
	data, err := SortByQueryParser(sortByString)
	fmt.Println(err)
	assert.Nil(t, err)
	assert.NotContains(t, "id", data)
	assert.NotContains(t, "createddate", data)
	assert.NotContains(t, "name", data)
}

func TestSortByParserAllowEmptyString(t *testing.T) {
	sortByString := ""
	data, err := SortByQueryParser(sortByString)
	fmt.Println(err)
	assert.Nil(t, err)
	assert.NotContains(t, "id", data)
	assert.NotContains(t, "createddate", data)
	assert.NotContains(t, "name", data)
}

func TestSortByParserAllowAnyCaseString(t *testing.T) {
	sortByString := `
	{	
		"ID":"AsC",
		"createdDate": "desc",
		"naMe": "Desc"
	}
	`
	data, err := SortByQueryParser(sortByString)
	fmt.Println(err)
	assert.Nil(t, err)
	assert.Equal(t, "asc", data["id"])
	assert.Equal(t, "desc", data["createddate"])
	assert.Equal(t, "desc", data["name"])
}

func TestSortByParserNotAllowEmptyValueForSoryByString(t *testing.T) {
	sortByString := `
	{	
		"ID":"",
		"createdDate": "desc",
		"naMe": "Desc"
	}
	`
	_, err := SortByQueryParser(sortByString)
	fmt.Println(err)
	assert.NotNil(t, err)
}
