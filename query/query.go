package query

import (
	"errors"
	"sort"
	"strings"

	"github.com/thomaspepio/hn-queries/constant"
	"github.com/thomaspepio/hn-queries/index"
	"github.com/thomaspepio/hn-queries/util"
)

// Query result : a single query with a count associated
type QueryResult struct {
	Query string `json:"query"`
	Count int    `json:"count"`
}

// CountURLs : counts URL occurences for the given couple datePrefix/keyType.
// Parameters datePrefix and keyType are assumed to be a match (e.g. datePrefix="2015" => keyType=util.Year)
func CountURLs(index *index.Index, datePrefix string, keyType util.KeyType) (int, error) {
	value, err := PerformSearch(index, datePrefix, keyType)

	if err != nil {
		return -1, err
	}

	countForDate := 0
	for _, urlCount := range value {
		countForDate += urlCount
	}

	return countForDate, nil
}

// FindTopNQueries : searches the top n queries for the given couple datePrefix/keyType.
// Parameters datePrefix and keyType are assumed to be a match (e.g. datePrefix="2015" => keyType=util.Year)
func FindTopNQueries(index *index.Index, datePrefix string, keyType util.KeyType, n int) ([]QueryResult, error) {
	value, err := PerformSearch(index, datePrefix, keyType)

	if err != nil {
		return nil, err
	}

	queriesForDate := make([]QueryResult, 0, len(value))
	for urlID, count := range value {
		url := index.IDstoURL[urlID]
		queriesForDate = append(queriesForDate, QueryResult{url, count})
	}

	sort.Slice(queriesForDate, func(i, j int) bool {
		return queriesForDate[i].Count > queriesForDate[j].Count
	})

	if n > len(queriesForDate) {
		return queriesForDate, nil
	} else {
		return queriesForDate[:n], nil
	}
}

// PerformSearch : perform a search on the index
func PerformSearch(index *index.Index, datePrefix string, keyType util.KeyType) (map[int]int, error) {
	var key int
	var keyError error
	switch keyType {
	case util.Year:
		key, keyError = util.YearKey(datePrefix)

		if keyError != nil {
			return nil, errors.New("Error during index search. Caused by : " + keyError.Error())
		}

	case util.Month:
		datePrefixSplitted := strings.Split(datePrefix, constant.Dash)
		key, keyError = util.MonthKey(datePrefixSplitted[0], datePrefixSplitted[1])

		if keyError != nil {
			return nil, errors.New("Error during index search. Caused by : " + keyError.Error())
		}

	case util.Day:
		datePrefixSplitted := strings.Split(datePrefix, constant.Dash)
		key, keyError = util.DayKey(datePrefixSplitted[0], datePrefixSplitted[1], datePrefixSplitted[2])

		if keyError != nil {
			return nil, errors.New("Error during index search. Caused by : " + keyError.Error())
		}
	}

	return index.Tree.Get(key), nil
}
