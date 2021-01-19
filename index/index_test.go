package index

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomaspepio/hn-queries/constant"
	"github.com/thomaspepio/hn-queries/parser"
)

func Test_AVLIndex_DeriveKeys(t *testing.T) {
	parsedQuery, _ := parser.ParseHNQuery(constant.CorrectLine)

	keys, _ := KeysFrom(parsedQuery)
	assert.Equal(t, 20150000000000, keys.Year, "Year should be 20150000000000")     // Our design requires that we can extraxt six different keys from each date, padding each possible key with zeroes when necessary.
	assert.Equal(t, 20150800000000, keys.Month, "Month should be 20150800000000")   // This poses a problem with the hours, minutes or seconds equal to 0.
	assert.Equal(t, 20150801000000, keys.Day, "Day should be 20150801240000")       // To remedy this : we choose to increment hours/seconds/minutes when represented in the key
	assert.Equal(t, 20150801010000, keys.Hour, "Hour should be 20150801010000")     // (e.g. there is an hour 24, a minute and a second 60).
	assert.Equal(t, 20150801010400, keys.Minute, "Second should be 20150801010400") // This allows us to keep the same order between dates and keys.
	assert.Equal(t, 20150801010444, keys.Second, "Minute should be 20150801010444")
}

func Test_AVLIndex_CannotParse(t *testing.T) {
	_, yearError := KeysFrom(nil)
	assert.Error(t, yearError, "You should not be able to parse from nil")
}

func Test_NewIndex_ShouldBeEmpty(t *testing.T) {
	index := EmptyIndex()
	assert.Equal(t, 0, index.Sequence, "Sequence should start at 0")
	assert.Empty(t, index.URLsToID, "No mapping from URL to ID should exist")
	assert.Empty(t, index.IDstoURL, "No mapping from ID to URL should exist")
	assert.Equal(t, 0, index.Tree.Height(), "Index tree should be empty")
}

func Test_AVLIndex_FromMultipleQueries(t *testing.T) {
	parsedQuery, _ := parser.ParseHNQuery(constant.CorrectLine)
	index := EmptyIndex()
	index.Add(parsedQuery)

	assert.Equal(t, len(index.URLsToID), 1, "One url should have been indexed")
	assert.Equal(t, len(index.IDstoURL), 1, "One url should have been indexed")
	assert.NotNil(t, index, "Index should not be nil")

	assert.NotNil(t, index.Tree.Get(20150000000000), "There should be a key for 20150000000000")
	assert.NotNil(t, index.Tree.Get(20150801010000), "There should be a key for 20150801010000")
	assert.NotNil(t, index.Tree.Get(20150800000000), "There should be a key for 20150800000000")
	assert.NotNil(t, index.Tree.Get(20150801000000), "There should be a key for 20150801000000")
	assert.NotNil(t, index.Tree.Get(20150801010400), "There should be a key for 20150801010400")
	//assert.NotNil(t, index.Tree.Get(20150801010444), "There should be a key for 20150801010444")
	assert.Equal(t, 1, len(index.Tree.Get(20150000000000)), "The key 20150000000000 should have seen one url")
	assert.Equal(t, 1, len(index.Tree.Get(20150800000000)), "The key 20150800000000 should have seen one url")
	assert.Equal(t, 1, len(index.Tree.Get(20150801000000)), "The key 20150801000000 should have seen one url")
	assert.Equal(t, 1, len(index.Tree.Get(20150801010000)), "The key 20150801010000 should have seen one url")
	assert.Equal(t, 1, len(index.Tree.Get(20150801010400)), "The key 20150801010400 should have seen one url")
	//assert.Equal(t, 1, len(index.Tree.Get(20150801010444)), "The key 20150801010444 should have seen one url")

	newURLAsString := "http://same-date-other-url"
	parsedQuery, _ = parser.ParseHNQuery(constant.DateAsString + constant.Tab + newURLAsString)
	index.Add(parsedQuery)
	assert.Equal(t, len(index.URLsToID), 2, "Two urls should have been indexed")
	assert.Equal(t, len(index.IDstoURL), 2, "Two urls should have been indexed")
	assert.Equal(t, 2, len(index.Tree.Get(20150000000000)), "The key 20150000000000 should have seen two urls")
	assert.Equal(t, 2, len(index.Tree.Get(20150800000000)), "The key 20150800000000 should have seen two urls")
	assert.Equal(t, 2, len(index.Tree.Get(20150801000000)), "The key 20150801000000 should have seen two urls")
	assert.Equal(t, 2, len(index.Tree.Get(20150801010000)), "The key 20150801010000 should have seen two url")
	assert.Equal(t, 2, len(index.Tree.Get(20150801010400)), "The key 20150801010400 should have seen two url")
	// assert.Equal(t, 2, len(index.Tree.Get(20150801010444)), "The key 20150801010444 should have seen two url")

	newDateAsString := "2021-01-01 00:03:43"
	newURLAsString = "http://other-url"
	parsedQuery, _ = parser.ParseHNQuery(newDateAsString + constant.Tab + newURLAsString)
	index.Add(parsedQuery)
	assert.Equal(t, len(index.URLsToID), 3, "Three urls should have been indexed")
	assert.Equal(t, len(index.IDstoURL), 3, "Three urls should have been indexed")
	assert.Equal(t, 1, len(index.Tree.Get(20210000000000)), "The key 20210000000000 should have seen one url")
	assert.Equal(t, 1, len(index.Tree.Get(20210000000000)), "The key 20210000000000 should have seen one url")
	assert.Equal(t, 1, len(index.Tree.Get(20210000000000)), "The key 20210000000000 should have seen one url")
	assert.Equal(t, 1, len(index.Tree.Get(20210101010000)), "The key 20210101010000 should have seen one url")
	assert.Equal(t, 1, len(index.Tree.Get(20210101010400)), "The key 20210101010400 should have seen one url")
	// assert.Equal(t, 1, len(index.Tree.Get(20210101010444)), "The key 20210101010444 should have seen one url")
}

func mapOf(key, val int) map[int]int {
	return map[int]int{key: val}
}
