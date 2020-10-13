package parsers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortByParser(t *testing.T) {
	sortByString := `id:asc;createdate:desc;name:desc`
	data, sortOrder, err := SortByQueryParser(sortByString)
	assert.Nil(t, err)
	assert.Equal(t, "asc", data["id"])
	assert.Equal(t, "desc", data["createdate"])
	assert.Equal(t, "desc", data["name"])
	assert.Equal(t, 3, len(sortOrder))
	assert.Equal(t, "id", sortOrder[0])
	assert.Equal(t, "createdate", sortOrder[1])
	assert.Equal(t, "name", sortOrder[2])
}

func TestSortByParserFailsOnUnknownFieldName(t *testing.T) {
	sortByString := "id:asc;createdate:desc;name:desc;invalid:desc"
	_, _, err := SortByQueryParser(sortByString)
	assert.NotNil(t, err)
	assert.Equal(t, "Unsupported value provided for sortBy", err.Error())
}

func TestSortByParserThrowErrorOnIllegalValue(t *testing.T) {
	sortByString := `"id":"randomValue";"createddate": "desc";"name": "desc"`
	_, _, err := SortByQueryParser(sortByString)
	assert.NotNil(t, err)
	assert.Equal(t, "Unsupported value provided for sortBy", err.Error())
}

func TestSortByParserAllowEmptyString(t *testing.T) {
	sortByString := ""
	data, sortOrder, err := SortByQueryParser(sortByString)
	fmt.Println(err)
	assert.Nil(t, err)
	assert.NotContains(t, "id", data)
	assert.NotContains(t, "createdate", data)
	assert.NotContains(t, "name", data)
	assert.Equal(t, 0, len(sortOrder))
}

func TestSortByParserAllowAnyCaseString(t *testing.T) {
	sortByString := "id:asc;creaTeDate:dEsc;naMe:desc"
	data, sortOrder, err := SortByQueryParser(sortByString)
	fmt.Println(err)
	assert.Nil(t, err)
	assert.Equal(t, "asc", data["id"])
	assert.Equal(t, "desc", data["createdate"])
	assert.Equal(t, "desc", data["name"])
	assert.Equal(t, 3, len(sortOrder))
	assert.Equal(t, "id", sortOrder[0])
	assert.Equal(t, "createdate", sortOrder[1])
	assert.Equal(t, "name", sortOrder[2])
}

func TestSortByParserNotAllowEmptyValueForSoryByString(t *testing.T) {
	sortByString := "id:;creaTeDate:dEsc;naMe:desc"
	_, _, err := SortByQueryParser(sortByString)
	assert.NotNil(t, err)
	assert.Equal(t, "Unsupported value provided for sortBy", err.Error())
}

func TestSortByParserAllowSemicolonAtTheEnd(t *testing.T) {
	sortByString := "id:desc;creaTeDate:dEsc;naMe:desc;"
	data, sortOrder, err := SortByQueryParser(sortByString)
	assert.Nil(t, err)
	assert.Equal(t, "desc", data["id"])
	assert.Equal(t, "desc", data["createdate"])
	assert.Equal(t, "desc", data["name"])
	assert.Equal(t, 3, len(sortOrder))
	assert.Equal(t, "id", sortOrder[0])
	assert.Equal(t, "createdate", sortOrder[1])
	assert.Equal(t, "name", sortOrder[2])
}

func TestPreserveOrderOfSortKeys(t *testing.T) {
	sortByString := "id:desc;creaTeDate:dEsc;naMe:desc;"
	_, sortOrder, err := SortByQueryParser(sortByString)
	assert.Nil(t, err)
	assert.Equal(t, "id", sortOrder[0])
	assert.Equal(t, "createdate", sortOrder[1])
	assert.Equal(t, "name", sortOrder[2])
}

func TestSortByParserAllowEmptyKeyValue(t *testing.T) {
	sortByString := "id:desc;creaTeDate:dEsc;naMe:desc;;"
	_, _, err := SortByQueryParser(sortByString)
	assert.NotNil(t, err)
	assert.Equal(t, "Unsupported value provided for sortBy", err.Error())
}
