package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thomaspepio/hn-queries/constant"
)

func Test_IdentifyKey_ShouldIdentifyYear(t *testing.T) {
	actualYear, _ := IdentifyKey("2021")
	assert.Equal(t, Year, actualYear, "Should have identified a year")
}

func Test_IdentifyKey_ShouldIdentifyMonth(t *testing.T) {
	actualMonth, _ := IdentifyKey("2021-01")
	assert.Equal(t, Month, actualMonth, "Should have identified a month")
}

func Test_IdentifyKey_ShouldIdentifyDay(t *testing.T) {
	actualDay, _ := IdentifyKey("2021-01-01")
	assert.Equal(t, Day, actualDay, "Should have identified a day")
}

func Test_IdentifyKey_ShouldIdentifyMinute(t *testing.T) {
	actualMinute, _ := IdentifyKey("2021-01-01 00:01")
	assert.Equal(t, Minute, actualMinute, "Should have identified a minute")
}

func Test_IdentifyKey_ShouldFail(t *testing.T) {
	_, err := IdentifyKey("YYYY")
	assert.Error(t, err, "YYYY is not a valid key")

	_, err = IdentifyKey("2021-MM")
	assert.Error(t, err, "2021-MM is not a valid key")

	_, err = IdentifyKey("2021-01-DD")
	assert.Error(t, err, "2021-01-DD is not a valid key")

	_, err = IdentifyKey("2021-01-01foobar")
	assert.Error(t, err, "2021-01-01foobar is not a valid key")
}

func Test_YearKey_ShouldSucceed(t *testing.T) {
	time, _ := time.Parse(constant.DateFormat, constant.DateAsString)
	key := YearKey(time)
	assert.Equal(t, 20150000000000, key, "Key should be 20150000000000")
}

func Test_MonthKey_ShouldSucceed(t *testing.T) {
	time, _ := time.Parse(constant.DateFormat, constant.DateAsString)
	key := MonthKey(time)
	assert.Equal(t, 20150800000000, key, "Key should be 20150800000000")
}

func Test_DayKey_ShouldSucceed(t *testing.T) {
	time, _ := time.Parse(constant.DateFormat, constant.DateAsString)
	key := DayKey(time)
	assert.Equal(t, 20150801000000, key, "Key should be 20150801000000")
}

func Test_HourKey_ShouldSucceed(t *testing.T) {
	time, _ := time.Parse(constant.DateFormat, constant.DateAsString)
	key := HourKey(time)
	assert.Equal(t, 20150801010000, key, "Key should be 20150801010000")
}

func Test_MinuteKey_ShouldSucceed(t *testing.T) {
	time, _ := time.Parse(constant.DateFormat, constant.DateAsString)
	key := MinuteKey(time)
	assert.Equal(t, 20150801010400, key, "Key should be 20150801010100")
}

func Test_SecondKey_ShouldSucceed(t *testing.T) {
	time, _ := time.Parse(constant.DateFormat, constant.DateAsString)
	key := SecondKey(time)
	assert.Equal(t, 20150801010444, key, "Key should be 20150801010101")
}
