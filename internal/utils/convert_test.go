package utils

import (
	"testing"
)

// convert string to int using strconv
func TestStrToInt(t *testing.T) {
	t.Run("wrong number", func(t *testing.T) {
		res := StrToInt("test")
		if res != -1 {
			t.Error("StrToInt(\"test\") should return -1")
		}
	})

	t.Run("good number", func(t *testing.T) {
		res := StrToInt("42")
		if res != 42 {
			t.Error("StrToInt(\"42\") should return 42")
		}
	})
}

// convert array of string to array of uuid
func TestStrToUUIDs(t *testing.T) {
	t.Run("wrong array of uuid", func(t *testing.T) {
		_, err := StrToUUIDs([]string{"test"})
		if err == nil {
			t.Error("StrToUUIDs([]string{\"test\"}) should return an error")
		}
	})

	t.Run("good array of uuid", func(t *testing.T) {
		_, err := StrToUUIDs([]string{"00000000-0000-0000-0000-000000000000"})
		if err != nil {
			t.Error("StrToUUIDs([]string{\"00000000-0000-0000-0000-000000000000\"}) should not return an error")
		}
	})

}

func TestToSnakeCase(t *testing.T) {
	t.Run("snake case", func(t *testing.T) {
		res := ToSnakeCase("snake_case")
		if res != "snake_case" {
			t.Error("ToSnakeCase(\"snake_case\") should return snake_case")
		}
	})

	t.Run("PascalCase", func(t *testing.T) {
		res := ToSnakeCase("PascalCase")
		if res != "pascal_case" {
			t.Error("ToSnakeCase(\"PascalCase\") should return pascal_case")
		}
	})

	t.Run("camelCase", func(t *testing.T) {
		res := ToSnakeCase("camelCase")
		if res != "camel_case" {
			t.Error("ToSnakeCase(\"camelCase\") should return camel_case")
		}
	})
}
