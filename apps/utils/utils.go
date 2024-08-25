package utils

import (
	"time"
)

func ParseDate(timestamp string) string {
	// Parse the timestamp into a time.Time object
	parsedTime, err := time.Parse("2006-01-02 15:04:05.000", timestamp)
	if err != nil {
		return ""
	}

	// Return only the date part as a string in the format "YYYY-MM-DD"
	return parsedTime.Format("2006-01-02")
}
