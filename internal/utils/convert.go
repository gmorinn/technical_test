package utils

import (
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/google/uuid"
)

// convert string to int using strconv
func StrToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}
	return i
}

// convert string to int64 using strconv
func StrToInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return -1
	}
	return i
}

// convert array of string to array of uuid
func StrToUUIDs(s []string) ([]uuid.UUID, error) {
	var res []uuid.UUID
	for _, id := range s {
		u, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil
}

// convert string to snake_case
func ToSnakeCase(s string) string {
	var res string
	for i, c := range s {
		if i == 0 {
			res += strings.ToLower(string(c))
		} else if c >= 'A' && c <= 'Z' {
			res += "_" + strings.ToLower(string(c))
		} else {
			res += string(c)
		}
	}
	return res
}

// convert string to PascalCase
func ToPascalCase(s string) string {
	var res string
	for i, c := range s {
		if i == 0 {
			res += strings.ToUpper(string(c))
		} else if c == '_' {
			continue
		} else if s[i-1] == '_' {
			res += strings.ToUpper(string(c))
		} else {
			res += string(c)
		}
	}
	return res
}

// convert string to UPPERCASE
func ToUpperCase(s string) string {
	return strings.ToUpper(s)
}

func ReadTemplate(filename string) string {
	// Read the file
	// Return the content
	f, err := os.ReadFile(filename)
	if err != nil {
		log.Println("error =>", err)
		return ""
	}
	return string(f)
}

func ReplacePlaceholders(s string, m map[string]string) string {
	for k, v := range m {
		s = strings.ReplaceAll(s, k, v)
	}
	return s
}

// func put upper case first letter
func PutUpperCaseFirstLetter(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// Helper function to lowercase the first letter of a string
func LowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

// add "[" and "]" to string, if mandatory add "[" and "!]" to string
func AddBrackets(s string, mandatory bool) string {
	if mandatory {
		return "[" + s + "!]"
	}
	return "[" + s + "]"
}
