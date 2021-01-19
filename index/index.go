package index

import (
	"errors"

	"github.com/thomaspepio/hn-queries/avltree"
	"github.com/thomaspepio/hn-queries/parser"
	"github.com/thomaspepio/hn-queries/util"
)

// Keys represents the different parts of an index key
// Input : 2021-01-17 11:22:33
// Key   : Year=20210000000000, Month=20210100000000, Day=20210117000000
//		   Hour=20210117110000, Minutes=20210117112200, Seconds=20210117112233
type IndexKeys struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
	Second int
}

// URLId : type alias for int
type URLId = int

// Index : a datastructure to deduplicate URLs and index them by year, year-month and year-month-day
type Index struct {
	Sequence int
	URLsToID map[string]URLId
	IDstoURL map[URLId]string
	Tree     *avltree.AVLTree
}

// EmptyIndex : creates an empty index
func EmptyIndex() *Index {
	almostEmptyTree := avltree.New(-1, make(map[int]int))
	return &Index{0, make(map[string]int), make(map[int]string), almostEmptyTree}
}

// Add : indexes a parsed query
func (index *Index) Add(parsedQuery *parser.ParsedQuery) error {
	keys, keysError := KeysFrom(parsedQuery)

	if keysError != nil {
		return keysError
	}

	//url := parsedQuery.URL.String()
	url := parsedQuery.URL
	urls := index.URLsToID
	ids := index.IDstoURL
	urlID, foundURL := urls[url]

	if !foundURL {
		sequence := index.Sequence
		urlID = sequence
		urls[url] = urlID
		ids[urlID] = url
	}
	index.Sequence++

	yearIndex := index.Tree.Get(keys.Year)
	if yearIndex == nil {
		index.Tree.Insert(keys.Year, initPairs(urlID))
	} else {
		yearIndex[urlID]++
		index.Tree.Update(keys.Year, yearIndex)
	}

	monthIndex := index.Tree.Get(keys.Month)
	if monthIndex == nil {
		index.Tree.Insert(keys.Month, initPairs(urlID))
	} else {
		monthIndex[urlID]++
		index.Tree.Update(keys.Month, monthIndex)
	}

	dayIndex := index.Tree.Get(keys.Day)
	if dayIndex == nil {
		index.Tree.Insert(keys.Day, initPairs(urlID))
	} else {
		dayIndex[urlID]++
		index.Tree.Update(keys.Day, dayIndex)
	}

	hourIndex := index.Tree.Get(keys.Hour)
	if hourIndex == nil {
		index.Tree.Insert(keys.Hour, initPairs(urlID))
	} else {
		hourIndex[urlID]++
		index.Tree.Update(keys.Hour, hourIndex)
	}

	minuteIndex := index.Tree.Get(keys.Minute)
	if minuteIndex == nil {
		index.Tree.Insert(keys.Minute, initPairs(urlID))
	} else {
		minuteIndex[urlID]++
		index.Tree.Update(keys.Minute, minuteIndex)
	}

	// secondIndex := index.Tree.Get(keys.Second)
	// if secondIndex == nil {
	// 	index.Tree.Insert(keys.Second, initPairs(urlID))
	// } else {
	// 	secondIndex[urlID]++
	// 	index.Tree.Update(keys.Second, secondIndex)
	// }

	return nil
}

// KeysFrom : parses a HN Query into a IndexKeys
func KeysFrom(parsedQuery *parser.ParsedQuery) (*IndexKeys, error) {
	if parsedQuery == nil {
		return nil, errors.New("Cannot extract index keys : no source to parse from")
	}

	time := parsedQuery.Time
	return &IndexKeys{
		Year:   util.YearKey(time),
		Month:  util.MonthKey(time),
		Day:    util.DayKey(time),
		Hour:   util.HourKey(time),
		Minute: util.MinuteKey(time),
		Second: util.SecondKey(time)}, nil
}

func initPairs(id int) map[int]int {
	return map[int]int{id: 1}
}
