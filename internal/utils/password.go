package utils

import (
	"errors"
	"math/rand"
)

func RandPasswordStringRunes() string {
	// need at least number, letter, uppercase, lowercase and 8 chars
	var uppercaseRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var lowercaseRunes = []rune("abcdefghijklmnopqrstuvwxyz")
	var numberRunes = []rune("0123456789")
	var specialRunes = []rune("!@#$%^&*()_+{}[]:;<>?/.,")
	b := make([]rune, 8)
	b[0] = uppercaseRunes[rand.Intn(len(uppercaseRunes))]
	b[1] = lowercaseRunes[rand.Intn(len(lowercaseRunes))]
	b[2] = numberRunes[rand.Intn(len(numberRunes))]
	b[3] = specialRunes[rand.Intn(len(specialRunes))]
	for i := 4; i < 8; i++ {
		b[i] = stringRunes[rand.Intn(len(stringRunes))]
	}
	return string(b)
}

func CheckPassword(password string) error {
	var uppercaseRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var lowercaseRunes = []rune("abcdefghijklmnopqrstuvwxyz")
	var numberRunes = []rune("0123456789")
	var specialRunes = []rune("!@#$%^&*()_+{}[]:;<>?/.,")
	var upper, lower, number, special bool

	for _, v := range password {
		for _, u := range uppercaseRunes {
			if v == u {
				upper = true
				break // Exit the loop once an uppercase letter is found.
			}
		}

		for _, l := range lowercaseRunes {
			if v == l {
				lower = true
				break // Exit the loop once a lowercase letter is found.
			}
		}

		for _, n := range numberRunes {
			if v == n {
				number = true
				break // Exit the loop once a number is found.
			}
		}

		for _, s := range specialRunes {
			if v == s {
				special = true
				break // Exit the loop once a special character is found.
			}
		}
	}

	if !upper {
		return errors.New("password must contain at least one uppercase letter")
	}

	if !lower {
		return errors.New("password must contain at least one lowercase letter")
	}

	if !number {
		return errors.New("password must contain at least one number")
	}

	if !special {
		return errors.New("password must contain at least one special character")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}
