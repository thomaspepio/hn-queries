package util

import (
	"errors"
	"regexp"
	"strconv"
	"time"
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

	// Format the API supports
	yearFormat   = "^[0-9]{4}$"
	monthFormat  = "^[0-9]{4}-[0-9]{2}$"
	dayFormat    = "^[0-9]{4}-[0-9]{2}-[0-9]{2}$"
	minuteFormat = "^[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}$"
)

var regexpYear = regexp.MustCompile(yearFormat)
var regexpMonth = regexp.MustCompile(monthFormat)
var regexpDay = regexp.MustCompile(dayFormat)
var regexpMinute = regexp.MustCompile(minuteFormat)

// IdentifyKey : associates a key string parameter to a supported API key type, or returns an error.
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

	if regexpMinute.MatchString(key) {
		return Minute, nil
	}

	return -1, errors.New("Could not identify key type from : " + key)
}

// YearKey : makes a search key for a whole year
func YearKey(time time.Time) int {
	return time.Year() * 10000000000
}

// MonthKey : makes a search key for a month in a year
func MonthKey(time time.Time) int {
	year := year(time.Year())
	month := month(time.Month())

	monthKey, _ := strconv.Atoi(year + month)
	return (monthKey * 100000000)
}

// DayKey : makes a search key for a specific day
func DayKey(time time.Time) int {
	year := year(time.Year())
	month := month(time.Month())
	day := day(time.Day())

	dayKey, _ := strconv.Atoi(year + month + day)
	return (dayKey * 1000000)
}

// HourKey : makes a search key for a specific hour
func HourKey(time time.Time) int {
	year := year(time.Year())
	month := month(time.Month())
	day := day(time.Day())
	hour := hour(time.Hour())

	hourKey, _ := strconv.Atoi(year + month + day + hour)
	return (hourKey * 10000)
}

// MinuteKey : makes a search key for a specific minute
func MinuteKey(time time.Time) int {
	year := year(time.Year())
	month := month(time.Month())
	day := day(time.Day())
	hour := hour(time.Hour())
	minute := minute(time.Minute())

	minuteKey, _ := strconv.Atoi(year + month + day + hour + minute)
	return (minuteKey * 100)
}

// SecondKey : makes a search key for a specific second
func SecondKey(time time.Time) int {
	year := year(time.Year())
	month := month(time.Month())
	day := day(time.Day())
	hour := hour(time.Hour())
	minute := minute(time.Minute())
	second := second(time.Second())

	secondKey, _ := strconv.Atoi(year + month + day + hour + minute + second)
	return secondKey
}

func year(year int) string {
	return strconv.Itoa(year)
}

func month(month time.Month) string {
	return pad(int(month))
}

func day(day int) string {
	return pad(day)
}

func hour(hour int) string {
	return padZeroableValue(hour)
}

func minute(minute int) string {
	return padZeroableValue(minute)
}

func second(second int) string {
	return padZeroableValue(second)
}

func pad(n int) string {
	if n < 10 {
		return "0" + strconv.Itoa(n)
	}

	return strconv.Itoa(n)
}

func padZeroableValue(n int) string {
	return pad(n + 1)
}
