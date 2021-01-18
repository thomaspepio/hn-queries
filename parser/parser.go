package parser

import (
	"errors"
	"strings"
	"time"

	"github.com/thomaspepio/hn-queries/constant"
)

// A ParsedQuery is a line parsed from the input file
type ParsedQuery struct {
	Time time.Time
	URL  string
}

// ParseHNQuery will parse a string with the following format : <YYYY-MM-DD HH:mm:SS><tab><url>
// An error is return when the string does not respect this format.
func ParseHNQuery(str string) (*ParsedQuery, error) {
	words := strings.Split(str, constant.Tab)

	var parsedQuery ParsedQuery
	if len(words) != 2 {
		return &parsedQuery, errors.New("Unable to parse line : " + str)
	}

	time, timeErr := time.Parse(constant.DateFormat, words[0])
	if timeErr != nil {
		return &parsedQuery, errors.New("Unable to parse date from : " + words[0])
	}

	url := words[1]

	parsedQuery = ParsedQuery{time, url}
	return &parsedQuery, nil
}
