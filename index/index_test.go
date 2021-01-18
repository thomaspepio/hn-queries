package index

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomaspepio/hn-queries/constant"
	"github.com/thomaspepio/hn-queries/parser"
)

func Test_AVLIndex_DeriveKeys(t *testing.T) {
	parsedQuery, _ := parser.ParseHNQuery(constant.CorrectLine)

	keysTriple, _ := KeysFrom(parsedQuery)
	assert.Equal(t, 20150000, keysTriple.Year, "Year should be 20150000")
	assert.Equal(t, 20150800, keysTriple.Month, "Month should be 20150800")
	assert.Equal(t, 20150801, keysTriple.Day, "Day should be 20150801")
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
	index := EmptyIndex()
	parsedQuery, _ := parser.ParseHNQuery(constant.CorrectLine)
	index.Add(parsedQuery)

	url := parsedQuery.URL.String()
	urlID := index.URLsToID[url]

	assert.Equal(t, len(index.URLsToID), 1, "One url should have been indexed")
	assert.Equal(t, len(index.IDstoURL), 1, "One url should have been indexed")
	assert.NotNil(t, index, "Index should not be nil")

	assert.NotNil(t, index.Tree.Get(20150000), "There should be a key for 20150000")
	assert.Equal(t, index.Tree.Get(20150000), mapOf(urlID, 1), "The key 20150000 should have seen one url")

	assert.NotNil(t, index.Tree.Get(20150800), "There should be a key for 20150800")
	assert.Equal(t, index.Tree.Get(20150800), mapOf(urlID, 1), "The key 20150800 should have seen one url")

	assert.NotNil(t, index.Tree.Get(20150801), "There should be a key for 20150801")
	assert.Equal(t, index.Tree.Get(20150801), mapOf(urlID, 1), "The key 20150801 should have seen one url")

	newURLAsString := "http://same-date-other-url"
	parsedQuery, _ = parser.ParseHNQuery(constant.DateAsString + constant.Tab + newURLAsString)
	index.Add(parsedQuery)
	assert.Equal(t, len(index.URLsToID), 2, "Two urls should have been indexed")
	assert.Equal(t, len(index.IDstoURL), 2, "Two urls should have been indexed")
	assert.Equal(t, 2, len(index.Tree.Get(20150000)), "The key 20150000 should have seen two urls")
	assert.Equal(t, 2, len(index.Tree.Get(20150800)), "The key 20150800 should have seen two urls")
	assert.Equal(t, 2, len(index.Tree.Get(20150801)), "The key 20150801 should have seen two urls")

	newDateAsString := "2021-01-01 00:03:43"
	newURLAsString = "http://other-url"
	parsedQuery, _ = parser.ParseHNQuery(newDateAsString + constant.Tab + newURLAsString)
	index.Add(parsedQuery)
	assert.Equal(t, len(index.URLsToID), 3, "Three urls should have been indexed")
	assert.Equal(t, len(index.IDstoURL), 3, "Three urls should have been indexed")
	assert.Equal(t, 2, len(index.Tree.Get(20150000)), "The key 20150000 should have seen two urls")
	assert.Equal(t, 2, len(index.Tree.Get(20150800)), "The key 20150800 should have seen two urls")
	assert.Equal(t, 2, len(index.Tree.Get(20150801)), "The key 20150801 should have seen two urls")
	assert.Equal(t, 1, len(index.Tree.Get(20210000)), "The key 20210000 should have seen one url")
	assert.Equal(t, 1, len(index.Tree.Get(20210100)), "The key 20210100 should have seen one url")
	assert.Equal(t, 1, len(index.Tree.Get(20210101)), "The key 20210101 should have seen one url")
}

func mapOf(key, val int) map[int]int {
	return map[int]int{key: val}
}
