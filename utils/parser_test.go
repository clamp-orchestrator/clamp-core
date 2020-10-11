package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFilter(t *testing.T) {
	filterString := `
	{
		"id": "asc",
		"createdDate": "desc",
		"name": "desc"
	}
	`
	data, err := ParseFilters(filterString)
	assert.Nil(t, err)
	assert.Equal(t, "asc", data["id"])
	assert.Equal(t, "desc", data["createdDate"])
	assert.Equal(t, "desc", data["name"])
}

func TestParseFilterAssignBlankToMissing(t *testing.T) {
	filterString := `
	{
		"createdDate": "desc",
		"name": "desc"
	}
	`
	data, err := ParseFilters(filterString)
	assert.Nil(t, err)
	assert.Equal(t, "", data["id"])
	assert.Equal(t, "desc", data["createdDate"])
	assert.Equal(t, "desc", data["name"])
}

func TestParseFilterThrowErrorOnIllegalValue(t *testing.T) {
	filterString := `
	{	
		"id":"randomValue",
		"createdDate": "desc",
		"name": "desc"
	}
	`
	_, err := ParseFilters(filterString)
	fmt.Println(err)
	assert.NotNil(t, err)
}

func TestParseFilterThrowErrorOnWrongJSON(t *testing.T) {
	filterString := `
	{	
		"id","randomValue",
		"createdDate": "desc",
		"name": "desc"
	}
	`
	_, err := ParseFilters(filterString)
	fmt.Println(err)
	assert.NotNil(t, err)
}
func TestParseFilterAllowEmptyJSON(t *testing.T) {
	filterString := "{}"
	data, err := ParseFilters(filterString)
	fmt.Println(err)
	assert.Nil(t, err)
	assert.Equal(t, "", data["id"])
	assert.Equal(t, "", data["createdDate"])
	assert.Equal(t, "", data["name"])
}

func TestParseFilterAllowEmptyString(t *testing.T) {
	filterString := ""
	data, err := ParseFilters(filterString)
	fmt.Println(err)
	assert.Nil(t, err)
	assert.Equal(t, "", data["id"])
	assert.Equal(t, "", data["createdDate"])
	assert.Equal(t, "", data["name"])
}
