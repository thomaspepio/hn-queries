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
	assert.Empty(t, index.Tree, "URL map should be empty")
	assert.Nil(t, index.Tree, "Index tree should be empty")
}

func Test_AVLIndex_FromMultipleQueries(t *testing.T) {
	parsedQuery, _ := parser.ParseHNQuery(constant.CorrectLine)
	url := parsedQuery.URL.String()

	index := EmptyIndex()
	index.Add(parsedQuery)
	urlId := index.URLs[url]

	assert.Equal(t, len(index.URLs), 1, "One url should have been indexed")
	assert.NotNil(t, index, "Index should not be nil")

	assert.NotNil(t, index.Tree.Get(20150000), "There should be a key for 20150000")
	assert.Equal(t, index.Tree.Get(20150000), mapOf(urlId, 1), "The key 20150000 should have counted "+url+"one time")

	assert.NotNil(t, index.Tree.Get(20150800), "There should be a key for 20150800")
	assert.Equal(t, index.Tree.Get(20150800), mapOf(urlId, 1), "The key 20150800 should have counted "+url+"one time")

	assert.NotNil(t, index.Tree.Get(20150801), "There should be a key for 20150801")
	assert.Equal(t, index.Tree.Get(20150801), mapOf(urlId, 1), "The key 20150801 should have counted "+url+"one time")

	index.Add(parsedQuery)
	assert.Equal(t, index.Tree.Get(20150000), mapOf(urlId, 2), "The key 20150000 should have counted "+url+"one time")
	assert.Equal(t, index.Tree.Get(20150800), mapOf(urlId, 2), "The key 20150800 should have counted "+url+"one time")
	assert.Equal(t, index.Tree.Get(20150801), mapOf(urlId, 2), "The key 20150801 should have counted "+url+"one time")

}

func mapOf(key, val int) map[int]int {
	return map[int]int{key: val}
}
