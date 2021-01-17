package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IdentifyKey_ShouldIdentifyYear(t *testing.T) {
	actualYear, _ := IdentifyKey("2021")
	assert.Equal(t, Year, actualYear, "Should have identified a year")
}

func Test_IdentifyKey_ShouldIdentifyMonth(t *testing.T) {
	actualMonth01, _ := IdentifyKey("2021-01")
	assert.Equal(t, Month, actualMonth01, "Should have identified a month")

	actualMonth12, _ := IdentifyKey("2021-12")
	assert.Equal(t, Month, actualMonth12, "Should have identified a month")
}

func Test_IdentifyKey_ShouldIdentifyDay(t *testing.T) {
	actualDay01, _ := IdentifyKey("2021-01-01")
	assert.Equal(t, Day, actualDay01, "Should have identified a day")

	actualDay27, _ := IdentifyKey("2021-01-27")
	assert.Equal(t, Day, actualDay27, "Should have identified a day")
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
	key, _ := YearKey("2015")
	assert.Equal(t, 20150000, key, "Key should be 20150000")
}

func Test_MonthKey_ShouldSucceed(t *testing.T) {
	key, _ := MonthKey("2015", "08")
	assert.Equal(t, 20150800, key, "Key should be 20150800")
}

func Test_DayKey_ShouldSucceed(t *testing.T) {
	key, _ := DayKey("2015", "08", "01")
	assert.Equal(t, 20150801, key, "YKey should be 20150801")
}
