package index

import (
	"errors"
	"strconv"

	"github.com/thomaspepio/hn-queries/avltree"
	"github.com/thomaspepio/hn-queries/parser"
	"github.com/thomaspepio/hn-queries/util"
)

// KeysTriple represents a 3-uple of keys
// With an input of : 2021-01-17,
// the key triple should hold Year=20210000, YearMonth=20210100, YearMonthDay=20210117
type KeysTriple struct {
	Year  int
	Month int
	Day   int
}

// URLId : type alias for int
type URLId = int

// Index : a datastructure to deduplicate URLs and index them by year, year-month and year-month-day
type Index struct {
	Sequence int
	URLs     map[string]URLId
	Tree     *avltree.AVLTree
}

// EmptyIndex : creates an empty index
func EmptyIndex() *Index {
	var emptyTree *avltree.AVLTree
	return &Index{0, make(map[string]int), emptyTree}
}

// Add : indexes a parsed query
func (index *Index) Add(parsedQuery *parser.ParsedQuery) error {
	keys, keysError := KeysFrom(parsedQuery)

	if keysError != nil {
		return keysError
	}

	url := parsedQuery.URL.String()
	urls := index.URLs
	urlID, foundURL := urls[url]

	if foundURL {
		yearPairs := index.Tree.Get(keys.Year)
		monthPairs := index.Tree.Get(keys.Month)
		dayPairs := index.Tree.Get(keys.Day)

		yearPairs[urlID]++
		monthPairs[urlID]++
		dayPairs[urlID]++
	} else {
		sequence := index.Sequence
		urlID := sequence
		sequence++
		urls[url] = urlID

		if index.Tree == nil {
			index.Tree = avltree.New(keys.Year, initPairs(urlID))
		} else {
			index.Tree.Insert(keys.Year, initPairs(urlID))
		}

		index.Tree.Insert(keys.Month, initPairs(urlID))
		index.Tree.Insert(keys.Day, initPairs(urlID))
	}

	return nil
}

// KeysFrom : parses a HN Query into a KeysTriple
func KeysFrom(parsedQuery *parser.ParsedQuery) (*KeysTriple, error) {
	if parsedQuery == nil {
		return nil, errors.New("Cannot extrat index keys : no source to parse from")
	}

	year, yearError := extractYearKey(parsedQuery)
	if yearError != nil {
		return nil, yearError
	}

	month, monthError := extractMonthKey(parsedQuery)
	if monthError != nil {
		return nil, monthError
	}

	day, dayError := extractDayKey(parsedQuery)
	if dayError != nil {
		return nil, dayError
	}

	return &KeysTriple{year, month, day}, nil
}

func extractYearKey(parsedQuery *parser.ParsedQuery) (int, error) {
	return util.YearKey(strconv.Itoa(parsedQuery.Time.Year()))
}

func extractMonthKey(parsedQuery *parser.ParsedQuery) (int, error) {
	year := strconv.Itoa(parsedQuery.Time.Year())
	month := pad(int(parsedQuery.Time.Month()))
	return util.MonthKey(year, month)
}

func extractDayKey(parsedQuery *parser.ParsedQuery) (int, error) {
	year := strconv.Itoa(parsedQuery.Time.Year())
	month := pad(int(parsedQuery.Time.Month()))
	day := pad(parsedQuery.Time.Day())

	return util.DayKey(year, month, day)
}

func pad(n int) string {
	if n < 10 {
		return "0" + strconv.Itoa(n)
	}

	return strconv.Itoa(n)
}

func initPairs(id int) map[int]int {
	return map[int]int{id: 1}
}
