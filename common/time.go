package common

import "time"

func ParseTime(timeString string) (time.Time, error) {
	time1, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		return time1, err
	}
	return time1, nil
}
