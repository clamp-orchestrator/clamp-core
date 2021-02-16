package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortByParser(t *testing.T) {
	sortByQuery := `id:asc,created_at:desc,name:desc`
	sortUsing, err := ParseFromQuery(sortByQuery)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(sortUsing))
	assert.Equal(t, "id", sortUsing[0].Key)
	assert.Equal(t, "asc", sortUsing[0].Order)
	assert.Equal(t, "created_at", sortUsing[1].Key)
	assert.Equal(t, "desc", sortUsing[1].Order)
	assert.Equal(t, "name", sortUsing[2].Key)
	assert.Equal(t, "desc", sortUsing[2].Order)
}

func TestSortByParserFailsOnUnknownFieldName(t *testing.T) {
	sortByQuery := "id:asc,created_at:desc,name:desc,invalid:desc"
	_, err := ParseFromQuery(sortByQuery)

	assert.NotNil(t, err)
	assert.Equal(t, "unsupported value provided for sortBy query", err.Error())
}

func TestSortByParserThrowErrorOnIllegalValue(t *testing.T) {
	sortByQuery := `"id":"randomValue","created_at": "desc","name": "desc"`
	_, err := ParseFromQuery(sortByQuery)

	assert.NotNil(t, err)
	assert.Equal(t, "unsupported value provided for sortBy query", err.Error())
}

func TestSortByParserAllowEmptyString(t *testing.T) {
	sortByQuery := ""
	sortUsing, err := ParseFromQuery(sortByQuery)

	assert.Nil(t, err)
	assert.Equal(t, 0, len(sortUsing))
}

func TestSortByParserAllowAnyCaseString(t *testing.T) {
	sortByQuery := `id:aSc,created_at:DeSc,NAME:desc`
	sortUsing, err := ParseFromQuery(sortByQuery)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(sortUsing))
	assert.Equal(t, "id", sortUsing[0].Key)
	assert.Equal(t, "asc", sortUsing[0].Order)
	assert.Equal(t, "created_at", sortUsing[1].Key)
	assert.Equal(t, "desc", sortUsing[1].Order)
	assert.Equal(t, "name", sortUsing[2].Key)
	assert.Equal(t, "desc", sortUsing[2].Order)
}

func TestSortByParserNotAllowEmptyValueForSoryByString(t *testing.T) {
	sortByQuery := "id:,created_at:dEsc,naMe:desc"
	_, err := ParseFromQuery(sortByQuery)

	assert.NotNil(t, err)
	assert.Equal(t, "unsupported value provided for sortBy query", err.Error())
}

func TestSortByParserAllowCommaAtTheEnd(t *testing.T) {
	sortByQuery := `id:asc,created_at:desc,name:desc,`
	sortUsing, err := ParseFromQuery(sortByQuery)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(sortUsing))
	assert.Equal(t, "id", sortUsing[0].Key)
	assert.Equal(t, "asc", sortUsing[0].Order)
	assert.Equal(t, "created_at", sortUsing[1].Key)
	assert.Equal(t, "desc", sortUsing[1].Order)
	assert.Equal(t, "name", sortUsing[2].Key)
	assert.Equal(t, "desc", sortUsing[2].Order)
}

func TestPreserveOrderOfSortKeys(t *testing.T) {
	sortByQuery := `created_at:desc,id:asc,name:desc,`

	sortUsing, err := ParseFromQuery(sortByQuery)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(sortUsing))
	assert.Equal(t, "created_at", sortUsing[0].Key)
	assert.Equal(t, "desc", sortUsing[0].Order)
	assert.Equal(t, "id", sortUsing[1].Key)
	assert.Equal(t, "asc", sortUsing[1].Order)
	assert.Equal(t, "name", sortUsing[2].Key)
	assert.Equal(t, "desc", sortUsing[2].Order)
}

func TestSortByParserAllowEmptyKeyValue(t *testing.T) {
	sortByQuery := "id:desc,created_at:dEsc,naMe:desc,,"
	_, err := ParseFromQuery(sortByQuery)

	assert.NotNil(t, err)
	assert.Equal(t, "unsupported value provided for sortBy query", err.Error())
}

func TestParseFromQueryReturnErrForInvalidField(t *testing.T) {
	sortByQuery := `"createdate": "desc"`
	_, err := ParseFromQuery(sortByQuery)

	assert.NotNil(t, err)
	assert.Equal(t, "unsupported value provided for sortBy query", err.Error())
}
