package endpoint

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/thomaspepio/hn-queries/query"

	"github.com/gin-gonic/gin"
	"github.com/thomaspepio/hn-queries/index"
	"github.com/thomaspepio/hn-queries/util"
)

const (
	v1queries = "/1/queries"

	// The datePrefix URL parameter name
	datePrefixParam = "datePrefix"

	// The size query parameter
	sizeParam = "size"

	// URLs we support
	countQueriesURL   = v1queries + "/count/:" + datePrefixParam
	popularQueriesURL = v1queries + "/popular/:" + datePrefixParam
)

// Router : return the endpoints of the application
func Router(index *index.Index) *gin.Engine {
	router := gin.Default()

	router.GET(countQueriesURL, func(context *gin.Context) {
		datePrefix := context.Param(datePrefixParam)

		keyType, keyTypeError := util.IdentifyKey(datePrefix)
		if keyTypeError != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "incorrect datePrefix parameter : " + keyTypeError.Error()})
		}

		count, countError := query.CountURLs(index, datePrefix, keyType)
		if countError != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "error while computing URL count : " + countError.Error()})
		} else {
			context.JSON(http.StatusOK, gin.H{"count": count})
		}
	})

	router.GET(popularQueriesURL, func(context *gin.Context) {
		datePrefix := context.Param(datePrefixParam)
		size := context.Query(sizeParam)

		keyType, datePrefixTypeError := util.IdentifyKey(datePrefix)
		if datePrefixTypeError != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": datePrefixTypeError.Error()})
		}

		n, sizeError := CheckSize(size)
		if sizeError != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": sizeError.Error()})
		} else {
			topQueries, topQueriesError := query.FindTopNQueries(index, datePrefix, keyType, n)
			if topQueriesError != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"error": "error while computing top queries : " + topQueriesError.Error()})
			} else {
				toJson, err := TopQueriesResultToJSON(topQueries)

				if err != nil {
					context.JSON(http.StatusInternalServerError, gin.H{"error": "error while computing top queries : " + err.Error()})
				}

				context.JSON(http.StatusOK, toJson)
			}
		}
	})

	return router
}

// CheckSize : checks the validity of the size query paramter
func CheckSize(size string) (int, error) {
	n, convError := strconv.Atoi(size)
	if convError != nil {
		return -1, errors.New("Wrong size parameter : " + size)
	}

	return n, nil
}

// TopQueriesResultToJSON : converts a []query.QueryResult to a JSON string
func TopQueriesResultToJSON(queryResults []query.QueryResult) (string, error) {
	var sb strings.Builder

	sb.WriteString("{queries:[")

	for _, query := range queryResults[:len(queryResults)-1] {
		asJSON, err := QueryResultToJson(query)

		if err != nil {
			return "error", err
		}

		sb.WriteString(asJSON)
		sb.WriteString(",")
	}

	asJSON, err := QueryResultToJson(queryResults[len(queryResults)-1])
	if err != nil {
		return "error", err
	}

	sb.WriteString(asJSON)
	sb.WriteString("]}")

	return sb.String(), nil
}

// QueryResult: converts a query.QueryResult to a JSON string
func QueryResultToJson(query query.QueryResult) (string, error) {
	byteArray, err := json.Marshal(query)

	if err != nil {
		return "error", err
	}

	queryQuotesRemoved := strings.Replace(string(byteArray), "\"query\"", "query", 1)
	countQuotesRemoved := strings.Replace(queryQuotesRemoved, "\"count\"", "count", 1)
	return countQuotesRemoved, nil
}
