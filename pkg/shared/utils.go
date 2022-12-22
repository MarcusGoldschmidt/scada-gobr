package shared

import (
	"strconv"
	"time"
)

// ParseDuration parses a duration string
func ParseDuration(minutes string) (time.Duration, error) {
	periodInt, err := strconv.ParseInt(minutes, 10, 64)

	if err != nil {
		return 0, err
	}

	return time.Duration(int64(time.Minute) * periodInt), nil
}

// ParseTime parses a time string
// Format: 2006-01-02T15:04:05Z07:00
func ParseTime(timeString string) (time.Time, error) {
	parse, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		return parse, err
	}

	return parse, nil
}
