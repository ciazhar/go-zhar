package date

import (
	"github.com/jinzhu/now"
	"strconv"
	"strings"
	"time"
)

func GetDateDistance(date, month, year int) (time.Time, time.Time) {

	i := getMonth(month)
	t := time.Date(year, i, date, 0, 0, 0, 0, time.UTC)
	return now.New(t).BeginningOfMonth(), now.New(t).EndOfMonth()
}

func StringToDate(date string) (time.Time, error) {
	dateArr := strings.Split(date, "-")
	day, err := strconv.Atoi(dateArr[0])
	if err != nil {
		log.Error(err.Error())
	}
	month, err := strconv.Atoi(dateArr[1])
	if err != nil {
		log.Error(err.Error())
	}
	year, err := strconv.Atoi(dateArr[2])
	if err != nil {
		log.Error(err.Error())
	}

	i := getMonth(month)

	return time.Date(year, i, day, 0, 0, 0, 0, time.UTC), nil

}

func getMonth(month int) time.Month {
	var i time.Month
	switch month {
	case 1:
		i = time.January
	case 2:
		i = time.February
	case 3:
		i = time.March
	case 4:
		i = time.April
	case 5:
		i = time.May
	case 6:
		i = time.June
	case 7:
		i = time.July
	case 8:
		i = time.August
	case 9:
		i = time.September
	case 10:
		i = time.October
	case 11:
		i = time.November
	case 12:
		i = time.December
	default:
		i = time.January
	}
	return i
}
