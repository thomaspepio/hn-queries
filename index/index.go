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
	Year         int
	YearMonth    int
	YearMonthDay int
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
	_, foundURL := urls[url]

	if foundURL {
		yearPair := index.Tree.Get(keys.Year)
		yearMonthPair := index.Tree.Get(keys.YearMonth)
		yearMonthDayPair := index.Tree.Get(keys.YearMonthDay)

		yearPair.Count++
		yearMonthPair.Count++
		yearMonthDayPair.Count++
	} else {
		sequence := index.Sequence
		urlID := sequence
		sequence++
		urls[url] = urlID

		if index.Tree == nil {
			index.Tree = avltree.New(keys.Year, newPair(urlID))
		} else {
			index.Tree.Insert(keys.Year, newPair(urlID))
		}

		index.Tree.Insert(keys.YearMonth, newPair(urlID))
		index.Tree.Insert(keys.YearMonthDay, newPair(urlID))
	}

	return nil
}

// KeysFrom : parses a HN Query into a KeysTriple
func KeysFrom(parsedQuery *parser.ParsedQuery) (*KeysTriple, error) {
	if parsedQuery == nil {
		return nil, errors.New("Cannot extrat index keys : no source to parse from")
	}

	year := extractYearKey(parsedQuery)

	yearMonth, yearMonthError := extractYearMonthKey(parsedQuery)
	if yearMonthError != nil {
		return nil, yearMonthError
	}

	yearMonthDay, yearMonthDayError := extractYearMonthDayKey(parsedQuery)
	if yearMonthDayError != nil {
		return nil, yearMonthDayError
	}

	return &KeysTriple{year, yearMonth, yearMonthDay}, nil
}

func extractYearKey(parsedQuery *parser.ParsedQuery) int {
	return util.MakeYearKey(parsedQuery.Time.Year())
}

func extractYearMonthKey(parsedQuery *parser.ParsedQuery) (int, error) {
	year := parsedQuery.Time.Year()
	month := int(parsedQuery.Time.Month())
	return util.MakeYearMonthKey(year, month)
}

func extractYearMonthDayKey(parsedQuery *parser.ParsedQuery) (int, error) {
	year := parsedQuery.Time.Year()
	month := int(parsedQuery.Time.Month())
	day := parsedQuery.Time.Day()

	return util.MakeYearMonthDayKey(year, month, day)
}

func pad(n int) string {
	if n < 10 {
		return "0" + strconv.Itoa(n)
	}

	return strconv.Itoa(n)
}

func newPair(id int) avltree.CountingPair {
	return avltree.CountingPair{id, 1}
}
