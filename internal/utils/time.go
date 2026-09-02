package utils

import "time"

// convert time.time to string
func TimeToString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// convert dd/mm/yyyy to time.Time
func StringToTime(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}
