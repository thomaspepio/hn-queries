package index

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomaspepio/hn-queries/avltree"
	"github.com/thomaspepio/hn-queries/constant"
	"github.com/thomaspepio/hn-queries/parser"
)

func Test_AVLIndex_DeriveKeys(t *testing.T) {
	parsedQuery, _ := parser.ParseHNQuery(constant.CorrectLine)

	keysTriple, _ := KeysFrom(parsedQuery)
	assert.Equal(t, 20150000, keysTriple.Year, "Year should be 20150000")
	assert.Equal(t, 20150800, keysTriple.YearMonth, "Year-Month should be 20150800")
	assert.Equal(t, 20150801, keysTriple.YearMonthDay, "Year should be 20150801")
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
	assert.NotNil(t, index.Tree.Get(20150000), "No key for 20150000")
	assert.Equal(t, *index.Tree.Get(20150000), avltree.CountingPair{urlId, 1}, "The key 20150000 should have counted "+url+"one time")
	assert.Equal(t, *index.Tree.Get(20150800), avltree.CountingPair{urlId, 1}, "The key 20150800 should have counted "+url+"one time")
	assert.Equal(t, *index.Tree.Get(20150801), avltree.CountingPair{urlId, 1}, "The key 20150801 should have counted "+url+"one time")

	index.Add(parsedQuery)
	assert.Equal(t, *index.Tree.Get(20150000), avltree.CountingPair{urlId, 2}, "The key 20150000 should have counted "+url+"one time")
	assert.Equal(t, *index.Tree.Get(20150800), avltree.CountingPair{urlId, 2}, "The key 20150800 should have counted "+url+"one time")
	assert.Equal(t, *index.Tree.Get(20150801), avltree.CountingPair{urlId, 2}, "The key 20150801 should have counted "+url+"one time")

}

func pairOf(val int) avltree.CountingPair {
	return avltree.CountingPair{val, 0}
}
