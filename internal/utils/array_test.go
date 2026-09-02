package utils

import "testing"

// check if the string is in the array
func TestInArray(t *testing.T) {
	t.Run("false case", func(t *testing.T) {
		res := InArray("test", []string{"test1", "test2"})
		if res == true {
			t.Error("InArray(\"test\", []string{\"test1\", \"test2\"}) should return false")
		}
	})

	t.Run("true case", func(t *testing.T) {
		res := InArray("test", []string{"test1", "test2", "test"})
		if res == false {
			t.Error("InArray(\"test\", []string{\"test1\", \"test2\", \"test\"}) should return true")
		}
	})
}
