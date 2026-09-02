package utils

import "strings"

func DecodeURIComponent(s string) (string, error) {
	// remove the %20 from the string
	res := strings.ReplaceAll(s, "%20", " ")
	return res, nil
}
