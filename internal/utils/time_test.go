package utils

import (
	"testing"
	"time"
)

// Testing TimeToString function
func TestTimeToString(t *testing.T) {
	now := time.Now()
	str := TimeToString(now)
	expected := now.Format("2006-01-02 15:04:05")
	if str != expected {
		t.Errorf("Expected %s, got %s", expected, str)
	}
}
