package query

import (
	"testing"

	"github.com/thomaspepio/hn-queries/util"

	"github.com/stretchr/testify/assert"
	"github.com/thomaspepio/hn-queries/constant"
	"github.com/thomaspepio/hn-queries/index"
	"github.com/thomaspepio/hn-queries/parser"
)

func Test_SearchIndex_SearchYeah_ShouldSucceed(t *testing.T) {
	index := index.EmptyIndex()

	parsedQuery, _ := parser.ParseHNQuery(constant.CorrectLine)
	index.Add(parsedQuery)

	pairs, _ := PerformSearch(index, "2015", util.Year)
	assert.NotNil(t, pairs, "Indexed content should have been found")
	assert.Equal(t, 1, len(pairs), "There should be only one indexed value for the query")
}

func Test_SearchIndex_SearchYeah_ShouldFail(t *testing.T) {
	index := index.EmptyIndex()

	parsedQuery, _ := parser.ParseHNQuery(constant.CorrectLine)
	index.Add(parsedQuery)

	_, err := PerformSearch(index, "YYYY", util.Year)
	assert.Error(t, err, "YYYY is not a valid API parameter")
}

func Test_SearchIndex_SearchMonth_ShouldSucceed(t *testing.T) {
	index := index.EmptyIndex()

	parsedQuery, _ := parser.ParseHNQuery(constant.CorrectLine)
	index.Add(parsedQuery)

	pairs, _ := PerformSearch(index, "2015-08", util.Month)
	assert.NotNil(t, pairs, "Indexed content should have been found")
	assert.Equal(t, 1, len(pairs), "There should be only one indexed value for the query")
}

func Test_SearchIndex_SearchMonth_ShouldFail(t *testing.T) {
	index := index.EmptyIndex()

	parsedQuery, _ := parser.ParseHNQuery(constant.CorrectLine)
	index.Add(parsedQuery)

	_, err := PerformSearch(index, "2015-MM", util.Year)
	assert.Error(t, err, "2015-MM is not a valid API parameter")
}

func Test_SearchIndex_SearchDay_ShouldSucceed(t *testing.T) {
	index := index.EmptyIndex()

	parsedQuery, _ := parser.ParseHNQuery(constant.CorrectLine)
	index.Add(parsedQuery)

	pairs, _ := PerformSearch(index, "2015-08-01", util.Day)
	assert.NotNil(t, pairs, "Indexed content should have been found")
	assert.Equal(t, 1, len(pairs), "There should be only one indexed value for the query")
}

func Test_SearchIndex_SearchDay_ShouldFail(t *testing.T) {
	index := index.EmptyIndex()

	parsedQuery, _ := parser.ParseHNQuery(constant.CorrectLine)
	index.Add(parsedQuery)

	_, err := PerformSearch(index, "2015-08-DD", util.Year)
	assert.Error(t, err, "2015-08-DD is not a valid API parameter")
}

func Test_ComputeKeyCount(t *testing.T) {
	parsedQuery, _ := parser.ParseHNQuery(constant.CorrectLine)
	index := index.EmptyIndex()
	index.Add(parsedQuery)

	value, _ := CountURLs(index, "2015", util.Year)
	assert.Equal(t, 1, value, "There should be one search for 2015")

	index.Add(parsedQuery)
	value, _ = CountURLs(index, "2015", util.Year)
	assert.Equal(t, 2, value, "There should be two searches for 2015")
}

func Test_TopNQueries(t *testing.T) {
	index := index.EmptyIndex()

	// Indexing the correct line from constants
	parsedQuery, _ := parser.ParseHNQuery(constant.CorrectLine)
	index.Add(parsedQuery)

	value, _ := FindTopNQueries(index, "2015", util.Year, 1)
	assert.Equal(t, 1, len(value), "There should be one top query since we indexed this url only once")

	value, _ = FindTopNQueries(index, "2015-08", util.Month, 1)
	assert.Equal(t, 1, len(value), "There should be one top query since we indexed this url only once")

	value, _ = FindTopNQueries(index, "2015-08-01", util.Day, 1)
	assert.Equal(t, 1, len(value), "There should be one top query since we indexed this url only once")

	// Indexing a bunch of different urls, at the same date than the previous one
	parsedQuery2, _ := parser.ParseHNQuery(constant.DateAsString + constant.Tab + "http://an-other-url")
	parsedQuery3, _ := parser.ParseHNQuery(constant.DateAsString + constant.Tab + "http://an-other-url2")
	parsedQuery4, _ := parser.ParseHNQuery(constant.DateAsString + constant.Tab + "http://an-other-url3")
	index.Add(parsedQuery2)
	index.Add(parsedQuery3)
	index.Add(parsedQuery4)

	value, _ = FindTopNQueries(index, "2015", util.Year, 4)
	assert.Equal(t, 4, len(value), "There should be four top query since we indexed four urls now")

	value, _ = FindTopNQueries(index, "2015-08", util.Month, 4)
	assert.Equal(t, 4, len(value), "There should be four top query since we indexed four urls now")

	value, _ = FindTopNQueries(index, "2015-08-01", util.Day, 4)
	assert.Equal(t, 4, len(value), "There should be four top query since we indexed four urls now")

	// Indexing a fifth URL, at a different time
	newDateAsString := "2021-01-01 00:03:43"
	newURLAsString := "http://and-again-an-other-url"
	parsedQuery, _ = parser.ParseHNQuery(newDateAsString + constant.Tab + newURLAsString)
	index.Add(parsedQuery)

	value, _ = FindTopNQueries(index, "2021", util.Year, 2)
	assert.Equal(t, 1, len(value), "There should be one top query since we indexed this url only once")

	value, _ = FindTopNQueries(index, "2021-01", util.Month, 2)
	assert.Equal(t, 1, len(value), "There should be one top query since we indexed this url only once")

	value, _ = FindTopNQueries(index, "2021-01-01", util.Day, 2)
	assert.Equal(t, 1, len(value), "There should be one top query since we indexed this url only once")
}
