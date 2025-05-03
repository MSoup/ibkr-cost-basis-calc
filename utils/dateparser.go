package utils

import (
	"time"
)

func ParseDate(dateStr string) (time.Time, error) {
	dateTime, err := time.Parse("2006-01-02, 15:04:05", dateStr)
	if err != nil {
		// Try alternative format if the first one fails
		dateTime, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return time.Time{}, err
		}
	}

	return dateTime, nil
}
