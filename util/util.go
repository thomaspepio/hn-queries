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
	Hour
	Minute
	Second

	yearFormat  = "^[0-9]{4}$"
	monthFormat = "^[0-9]{4}-[0-9]{1,2}$"
	dayFormat   = "^[0-9]{4}-[0-9]{1,2}-[0-9]{1,2}$"
)

var regexpYear = regexp.MustCompile(yearFormat)
var regexpMonth = regexp.MustCompile(monthFormat)
var regexpDay = regexp.MustCompile(dayFormat)

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
func YearKey(yearStr string) int {
	yearKey, _ := strconv.Atoi(yearStr)
	return yearKey * 10000000000
}

// MonthKey : makes a search key for a month in a year
func MonthKey(year, month string) int {
	yearMonthStr := year + month
	monthKey, _ := strconv.Atoi(yearMonthStr)
	return (monthKey * 100000000)
}

// DayKey : makes a search key for a specific day
func DayKey(year, month, day string) int {
	dayKey, _ := strconv.Atoi(year + month + day)
	return (dayKey * 1000000)
}

// HourKey : makes a search key for a specific hour
func HourKey(year, month, day, hour string) int {
	hourKey, _ := strconv.Atoi(year + month + day + hour)
	return (hourKey * 10000)
}

// MinuteKey : makes a search key for a specific minute
func MinuteKey(year, month, day, hour, minute string) int {
	minuteKey, _ := strconv.Atoi(year + month + day + hour + minute)
	return (minuteKey * 100)
}

// SecondKey : makes a search key for a specific second
func SecondKey(year, month, day, hour, minute, second string) int {
	secondKey, _ := strconv.Atoi(year + month + day + hour + minute + second)
	return secondKey
}
