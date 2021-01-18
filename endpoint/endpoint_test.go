package endpoint

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomaspepio/hn-queries/query"
)

func Test_Size_AnyNumber_ShouldBeAccepted(t *testing.T) {
	size := rand.Intn(100)
	n, err := CheckSize(strconv.Itoa(size))
	assert.NotNil(t, n, "Any number is an acceptable size parameter")
	assert.Nil(t, err, "Any number is an acceptable size parameter")
}

func Test_Size_OtherThanNumber_ShouldNotBeAccepted(t *testing.T) {
	_, err := CheckSize("foo")
	assert.Error(t, err, "Anything that is not a number is not a valid size parameter")
}

func Test_TopQueryResult_ToJSON(t *testing.T) {
	queryResult := query.QueryResult{"foo", 1}
	asJson, _ := QueryResultToJson(queryResult)
	assert.Equal(t, "{query:\"foo\",count:1}", asJson, "QueryResult not JSON encoded properly")
}

func Test_TopQueryResultArray_ToJSON(t *testing.T) {
	queryResult1 := query.QueryResult{"foo", 1}
	queryResult2 := query.QueryResult{"bar", 2}
	queries := []query.QueryResult{queryResult1, queryResult2}

	asJson, _ := TopQueriesResultToJSON(queries)
	assert.Equal(t, "{queries:[{query:\"foo\",count:1},{query:\"bar\",count:2}]}", asJson, "TopQueriesResult not JSON encoded propertly")
}
