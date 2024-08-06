package sqlite

import (
	"strings"
	"time"
)

func parseDbTimeString(i *string) (time.Time, error) {
	if i == nil {
		return time.Time{}, nil
	}

	if strings.Contains(*i, "T") {
		if strings.Contains(*i, "Z") {
			return time.Parse("2006-01-02T15:04:05Z", *i)
		}

		return time.Parse("2006-01-02T15:04:05-07:00", *i)
	}

	if strings.Contains(*i, "Z") {
		return time.Parse("2006-01-02 15:04:05Z", *i)
	}

	return time.Parse("2006-01-02 15:04:05-07:00", *i)
}
