package fizzbuzz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate_Classic(t *testing.T) {
	result := Generate(3, 5, 16, "fizz", "buzz")
	expected := []string{
		"1", "2", "fizz", "4", "buzz",
		"fizz", "7", "8", "fizz", "buzz",
		"11", "fizz", "13", "14", "fizzbuzz",
		"16",
	}
	assert.Equal(t, expected, result)
}

func TestGenerate_Limit1(t *testing.T) {
	result := Generate(3, 5, 1, "fizz", "buzz")
	assert.Equal(t, []string{"1"}, result)
}

func TestGenerate_MultiplesOfBoth(t *testing.T) {
	result := Generate(2, 3, 6, "foo", "bar")
	expected := []string{"1", "foo", "bar", "foo", "5", "foobar"}
	assert.Equal(t, expected, result)
}
