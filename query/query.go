package query

import (
	"errors"
	"sort"
	"time"

	"github.com/thomaspepio/hn-queries/avltree"
	"github.com/thomaspepio/hn-queries/constant"
	"github.com/thomaspepio/hn-queries/index"
	"github.com/thomaspepio/hn-queries/util"
)

const (
	yearFormat   = "2006"
	monthFormat  = "2006-01"
	dayFormat    = "2006-01-02"
	hourFormat   = "2006-01-02 15"
	minuteFormat = "2006-01-02 15:04"
	secondFormat = "2006-01-02 15:04:05"
)

// QueryResult : a single query with a count associated
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
	}

	return queriesForDate[:n], nil
}

// PerformSearch : perform a search on the index
func PerformSearch(index *index.Index, datePrefix string, keyType util.KeyType) (map[int]int, error) {
	var key int

	switch keyType {
	case util.Year:
		datePrefixAsTime, _ := time.Parse(yearFormat, datePrefix) // TODO tester le cas erreur
		key = util.YearKey(datePrefixAsTime)
		return index.Tree.Get(key), nil

	case util.Month:
		datePrefixAsTime, _ := time.Parse(monthFormat, datePrefix) // TODO tester le cas erreur
		key = util.MonthKey(datePrefixAsTime)
		return index.Tree.Get(key), nil

	case util.Day:
		datePrefixAsTime, _ := time.Parse(dayFormat, datePrefix) // TODO tester le cas erreur
		key = util.DayKey(datePrefixAsTime)
		return index.Tree.Get(key), nil

	case util.Minute:
		lower, _ := time.Parse(minuteFormat, datePrefix) // TODO tester le cas erreur
		lowerKey := util.MinuteKey(lower)
		higherKey := lowerKey + 60

		treeBetween := index.Tree.Between(lowerKey, higherKey)
		return mostGeneralTreeValue(treeBetween), nil
	}

	return nil, errors.New("No key was extracted. This is an error")
}

func splitter(r rune) bool {
	return r == constant.DashRune || r == constant.SpaceRune || r == constant.ColonRune
}

func mostGeneralTreeValue(tree *avltree.AVLTree) map[int]int {
	if tree == nil {
		return nil
	}

	var result map[int]int

	if tree.Left == nil {
		result = tree.Values
	} else {
		result = mostGeneralTreeValue(tree.Left)
	}

	return result
}
