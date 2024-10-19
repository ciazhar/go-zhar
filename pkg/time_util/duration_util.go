package timeutil

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

// ParseTimeDuration parses a string with time units (h, m, s) and returns time.Duration
func ParseTimeDuration(input string) (time.Duration, error) {
	// Define regex pattern to match time units (e.g., 12h, 2m, 4s)
	re := regexp.MustCompile(`(\d+)([hms])`)
	matches := re.FindAllStringSubmatch(input, -1)

	if matches == nil {
		return 0, errors.New("invalid time format")
	}

	var totalDuration time.Duration
	for _, match := range matches {
		// match[1] is the number part, match[2] is the unit part (h/m/s)
		value, err := strconv.Atoi(match[1])
		if err != nil {
			return 0, err
		}

		switch match[2] {
		case "h":
			totalDuration += time.Duration(value) * time.Hour
		case "m":
			totalDuration += time.Duration(value) * time.Minute
		case "s":
			totalDuration += time.Duration(value) * time.Second
		default:
			return 0, errors.New("unknown time unit")
		}
	}

	return totalDuration, nil
}