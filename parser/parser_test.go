package parser

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thomaspepio/hn-queries/constant"
)

func Test_ValidHNQuery_ShouldBeParsed(t *testing.T) {
	timeParsed, _ := time.Parse(constant.DateFormat, constant.DateAsString)
	urlParsed, _ := url.Parse(constant.URLAsString)

	expectedQuery := ParsedQuery{timeParsed, *urlParsed}
	parsedQuery, _ := ParseHNQuery(constant.CorrectLine)

	assert.Equal(t,
		expectedQuery,
		*parsedQuery,
		"A valid HN query should be parsed")
}

func Test_InvalidLine_ShouldReturnError(t *testing.T) {
	invalidLine := "absolutely not a correct line"
	_, err := ParseHNQuery(invalidLine)
	assert.EqualError(t,
		err,
		"Unable to parse line : "+invalidLine,
		"An invalid line should not be parsed")
}

func Test_InvalidDate_ShouldReturnError(t *testing.T) {
	_, err := ParseHNQuery("not-a-date" + constant.Tab + constant.URLAsString)
	assert.EqualError(t,
		err,
		"Unable to parse date from : not-a-date",
		"An invalid date should not be parsed")
}
