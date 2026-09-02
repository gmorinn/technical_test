package utils

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
)

// NullS return sql.NullString from string
func TestNullS(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		res := NullS("")
		if res.String != "" || res.Valid != false {
			t.Error("NullS(\"\") should return sql.NullString{}")
		}
	})

	t.Run("not empty string", func(t *testing.T) {
		res := NullS("test")
		if res.String != "test" || res.Valid != true {
			t.Error("NullS(\"test\") should return sql.NullString{String: \"test\", Valid: true}")
		}
	})
}

// NullT return sql.NullTime from time.Time
func TestNullT(t *testing.T) {
	t.Run("zero time", func(t *testing.T) {
		res := NullT(time.Time{})
		if res.Valid != false {
			t.Error("NullT(time.Time{}) should return sql.NullTime{}")
		}
	})

	t.Run("not zero time", func(t *testing.T) {
		res := NullT(time.Now())
		if res.Valid != true {
			t.Error("NullT(time.Now()) should return sql.NullTime{Time: time.Now(), Valid: true}")
		}
	})
}

// NullB return sql.NullBool from bool
func TestNullB(t *testing.T) {
	t.Run("false bool", func(t *testing.T) {
		res := NullB(false)
		if res.Bool != false || res.Valid != true {
			t.Error("NullB(false) should return sql.NullBool{Bool: false, Valid: true}")
		}
	})

	t.Run("true bool", func(t *testing.T) {
		res := NullB(true)
		if res.Bool != true || res.Valid != true {
			t.Error("NullB(true) should return sql.NullBool{Bool: true, Valid: true}")
		}
	})
}

// NullU return empty string instedof 00000000-0000-0000-0000-000000000000
func TestNullU(t *testing.T) {
	t.Run("empty uuid", func(t *testing.T) {
		res := NullU(uuid.Nil)
		if res != "" {
			t.Error("NullU(uuid.Nil) should return \"\"")
		}
	})

	t.Run("not empty uuid", func(t *testing.T) {
		res := NullU(uuid.New())
		if res == "" {
			t.Error("NullU(uuid.New()) should not return \"\"")
		}
	})
}

// PointeurS return sql.NullString from string
func TestPointeurSToNullS(t *testing.T) {
	t.Run("nil string", func(t *testing.T) {
		res := PointeurSToNullS(nil)
		if res.String != "" || res.Valid != false {
			t.Error("PointeurSToNullS(nil) should return sql.NullString{}")
		}
	})

	t.Run("not nil string", func(t *testing.T) {
		var str = "test"
		res := PointeurSToNullS(&str)
		if res.String != "test" || res.Valid != true {
			t.Error("PointeurSToNullS(&str) should return sql.NullString{String: \"test\", Valid: true}")
		}
	})
}

// Sql.NullString to *string
func TestNullSToPointeurS(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		res := NullSToPointeurS(sql.NullString{})
		if res != nil {
			t.Error("NullSToPointeurS(sql.NullString{}) should return nil")
		}
	})

	t.Run("not empty string", func(t *testing.T) {
		res := NullSToPointeurS(sql.NullString{String: "test", Valid: true})
		if *res != "test" {
			t.Error("NullSToPointeurS(sql.NullString{String: \"test\", Valid: true}) should return \"test\"")
		}
	})
}

// Sql.NullInt32 to *int
func TestNullI32ToPointeurI(t *testing.T) {
	t.Run("empty int", func(t *testing.T) {
		res := NullI32ToPointeurI(sql.NullInt32{})
		if res != nil {
			t.Error("NullI32ToPointeurI32(sql.NullInt32{}) should return nil")
		}
	})

	t.Run("not empty int", func(t *testing.T) {
		res := NullI32ToPointeurI(sql.NullInt32{Int32: 42, Valid: true})
		if *res != 42 {
			t.Error("NullI32ToPointeurI32(sql.NullInt32{Int32: 42, Valid: true}) should return 42")
		}
	})
}

// Sql.NullInt64 to *int64
func TestNullI64ToPointeurI64(t *testing.T) {
	t.Run("empty int64", func(t *testing.T) {
		res := NullI64ToPointeurI64(sql.NullInt64{})
		if res != nil {
			t.Error("NullI64ToPointeurI64(sql.NullInt64{}) should return nil")
		}
	})

	t.Run("not empty int64", func(t *testing.T) {
		res := NullI64ToPointeurI64(sql.NullInt64{Int64: 42, Valid: true})
		if *res != 42 {
			t.Error("NullI64ToPointeurI64(sql.NullInt64{Int64: 42, Valid: true}) should return 42")
		}
	})
}

// check email format
func TestIsEmail(t *testing.T) {
	t.Run("wrong email", func(t *testing.T) {
		res := IsEmail("test")
		if res != false {
			t.Error("IsEmail(\"test\") should return false")
		}
	})

	t.Run("good email", func(t *testing.T) {
		res := IsEmail("test@gmail.com")
		if res != true {
			t.Error("IsEmail(\"test@gmail.com\") should return true")
		}
	})
}

// check uuid format
func TestIsUUID(t *testing.T) {
	t.Run("wrong uuid", func(t *testing.T) {
		res := IsUUID("test")
		if res != false {
			t.Error("IsUUID(\"test\") should return false")
		}
	})

	t.Run("good uuid", func(t *testing.T) {
		res := IsUUID("00000000-0000-0000-0000-000000000000")
		if res != true {
			t.Error("IsUUID(\"00000000-0000-0000-0000-000000000000\") should return true")
		}
	})
}

// Testing NullF, NullI64, NullI32, NullI32ToPointeurI32 functions
func TestNullF(t *testing.T) {
	nullFloat := NullF(1.2)
	if !nullFloat.Valid || nullFloat.Float64 != 1.2 {
		t.Error("NullFloat64 is not valid or does not contain the correct value")
	}
}

func TestNullI64(t *testing.T) {
	nullInt64 := NullI64(10)
	if !nullInt64.Valid || nullInt64.Int64 != 10 {
		t.Error("NullInt64 is not valid or does not contain the correct value")
	}
}

func TestNullI32(t *testing.T) {
	nullInt32 := NullI32(10)
	if !nullInt32.Valid || nullInt32.Int32 != 10 {
		t.Error("NullInt32 is not valid or does not contain the correct value")
	}
}

func TestNullI32ToPointeurI32(t *testing.T) {
	nullInt32 := NullI32(10)
	ptr := NullI32ToPointeurI32(nullInt32)
	if ptr == nil || *ptr != 10 {
		t.Error("Pointer is null or does not point to the correct value")
	}

	res := NullI32ToPointeurI32(sql.NullInt32{})
	if res != nil {
		t.Error("NullI32ToPointeurI(sql.NullInt32{}) should return nil")
	}

}
