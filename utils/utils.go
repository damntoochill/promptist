package utils

import (
	"database/sql"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

func Random(min int, max int) int {
	return rand.Intn(max-min) + min
}

// NewNullString todo
func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// NewNullInt todo
func NewNullInt(s string) sql.NullInt64 {
	if len(s) == 0 {
		return sql.NullInt64{}
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		fmt.Printf("%d of type %T", n, n)
	}
	return sql.NullInt64{
		Int64: n,
		Valid: true,
	}
}

// NewNullFloat todo
func NewNullFloat(s string) sql.NullFloat64 {
	if len(s) == 0 {
		return sql.NullFloat64{}
	}
	n, err := strconv.ParseFloat(s, 8)
	if err == nil {
		fmt.Println(n, err, reflect.TypeOf(n))
	}
	return sql.NullFloat64{
		Float64: n,
		Valid:   true,
	}
}

// NewNullDate todo
func NewNullDate(s string) sql.NullTime {
	if len(s) == 0 {
		return sql.NullTime{}
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		fmt.Println(err)
		return sql.NullTime{}
	}
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}
