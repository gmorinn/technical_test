package utils

import (
	"strings"
	"testing"
	"unicode"
)

// Testing RandStringRunes function
func TestRandStringRunes(t *testing.T) {
	str := RandStringRunes(5)
	if len(str) != 5 {
		t.Error("String length does not match the expected length")
	}
}

// Testing RandNumberRunes function
func TestRandNumberRunes(t *testing.T) {
	num := RandNumberRunes(5)
	if len(num) != 5 {
		t.Error("Number string length does not match the expected length")
	}
}

// Testing RandomInt function
func TestRandomInt(t *testing.T) {
	randNum := RandomInt(5, 10)
	if randNum < 5 || randNum > 10 {
		t.Error("Random number is out of the expected range")
	}
}

// Testing RandomAttribut function
func TestRandomAttribut(t *testing.T) {
	arr := []string{"A", "B", "C"}
	randAttr := RandomAttribut(arr)
	if !strings.Contains("ABC", randAttr) {
		t.Error("Random attribute is not part of the input array")
	}
}

// Testing RandLetterRunes function
func TestRandLetterRunes(t *testing.T) {
	str := RandLetterRunes(5)
	if len(str) != 5 {
		t.Error("String length does not match the expected length")
	}

	for _, char := range str {
		if !unicode.IsLetter(rune(char)) {
			t.Error("String contains non-letter characters")
		}
	}
}
