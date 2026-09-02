package utils

import (
	"testing"
	"unicode"
)

func TestRandPasswordStringRunes(t *testing.T) {
	str := RandPasswordStringRunes()
	if len(str) != 8 {
		t.Error("String length does not match the expected length")
	}

	// Check if the string contains at least one uppercase letter
	containsUppercase := false
	for _, char := range str {
		if unicode.IsUpper(rune(char)) {
			containsUppercase = true
			break
		}
	}
	if !containsUppercase {
		t.Error("String does not contain at least one uppercase letter")
	}

	// Check if the string contains at least one lowercase letter
	containsLowercase := false
	for _, char := range str {
		if unicode.IsLower(rune(char)) {
			containsLowercase = true
			break
		}
	}
	if !containsLowercase {
		t.Error("String does not contain at least one lowercase letter")
	}

	// Check if the string contains at least one number
	containsNumber := false
	for _, char := range str {
		if unicode.IsNumber(rune(char)) {
			containsNumber = true
			break
		}
	}
	if !containsNumber {
		t.Error("String does not contain at least one number")
	}

	// Check if the string contains at least one special character.
	//
	// Both categories are accepted: RandPasswordStringRunes draws from
	// "!@#$%^&*()_+{}[]:;<>?/.," and five of those — $ ^ + < > — are Unicode
	// symbols (Sc/Sk/Sm) rather than punctuation. Testing IsPunct alone made
	// this test fail whenever the generator picked one of them.
	containsSpecial := false
	for _, char := range str {
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			containsSpecial = true
			break
		}
	}
	if !containsSpecial {
		t.Errorf("String %q does not contain at least one special character", str)
	}
}

func TestCheckPassword(t *testing.T) {
	// Test if the password is at least 8 characters long
	err := CheckPassword("Abcdefg1")
	if err == nil {
		t.Error("Password is not at least 8 characters long")
	}

	// Test if the password contains at least one uppercase letter
	err = CheckPassword("abcdefg1")
	if err == nil {
		t.Error("Password does not contain at least one uppercase letter")
	}

	// Test if the password contains at least one lowercase letter
	err = CheckPassword("ABCDEFG1")
	if err == nil {
		t.Error("Password does not contain at least one lowercase letter")
	}

	// Test if the password contains at least one number
	err = CheckPassword("Abcdefgh")
	if err == nil {
		t.Error("Password does not contain at least one number")
	}

	// Test if the password contains at least one special character
	err = CheckPassword("Abcdefg1")
	if err == nil {
		t.Error("Password does not contain at least one special character")
	}

	// Test if the password is good
	err = CheckPassword("Ebcdefg1!")
	if err != nil {
		t.Error("Password should be good")
	}
}
