package util

import (
	"errors"
	"regexp"
	"strconv"
)

type KeyType int

const (
	// Enumerations of possible key shapes
	Year KeyType = iota
	Month
	Day

	datePrefixParamFormatYear  = "^[0-9]{4}$"
	datePrefixParamFormatMonth = "^[0-9]{4}-[0-9]{1,2}$"
	datePrefixParamFormathDay  = "^[0-9]{4}-[0-9]{1,2}-[0-9]{1,2}$"
)

var regexpYear = regexp.MustCompile(datePrefixParamFormatYear)
var regexpMonth = regexp.MustCompile(datePrefixParamFormatMonth)
var regexpDay = regexp.MustCompile(datePrefixParamFormathDay)

func IdentifyKey(key string) (KeyType, error) {
	if regexpYear.MatchString(key) {
		return Year, nil
	}

	if regexpMonth.MatchString(key) {
		return Month, nil
	}

	if regexpDay.MatchString(key) {
		return Day, nil
	}

	return -1, errors.New("Could not identify key from : " + key)
}

// YearKey : makes a search key for a whole year
// Year = 2021 => Key = 20210000
func YearKey(yearStr string) (int, error) {
	year, err := strconv.Atoi(yearStr)

	if err != nil {
		return 0, err
	}

	return year * 10000, nil
}

// MonthKey : makes a search key for a month in a year
// Year = 2021, Month = 01 => Key = 20210100
func MonthKey(year, month string) (int, error) {
	yearMonthStr := year + month
	yearMonth, convError := strconv.Atoi(yearMonthStr)

	if convError != nil {
		return 0, convError
	}

	return (yearMonth * 100), nil
}

// DayKey : makes a search key for a specific day
// Year = 2021, Month = 01, Day = 01 => Key = 20210101
func DayKey(year, month, day string) (int, error) {
	yearMonthDay, convError := strconv.Atoi(year + month + day)

	if convError != nil {
		return 0, convError
	}

	return yearMonthDay, nil
}
