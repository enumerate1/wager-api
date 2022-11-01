package database

import (
	"time"

	"github.com/jackc/pgtype"
)

// Text converts a Go string to pgtype.Text.
func Text(v string) pgtype.Text {
	return pgtype.Text{String: v, Status: pgtype.Present}
}

// Int4 converts a Go int32 to pgtype.Int4.
func Int4(v int32) pgtype.Int4 {
	return pgtype.Int4{Int: v, Status: pgtype.Present}
}

// Int8 converts a Go int64 to pgtype.Int8.
func Int8(v int64) pgtype.Int8 {
	return pgtype.Int8{Int: v, Status: pgtype.Present}
}

// Float4 converts a Go float32 to pgtype.Float4.
func Float4(v float32) pgtype.Float4 {
	return pgtype.Float4{Float: v, Status: pgtype.Present}
}

// Bool converts a Go bool to pgtype.Bool.
func Bool(v bool) pgtype.Bool {
	return pgtype.Bool{Bool: v, Status: pgtype.Present}
}

// Timestamptz converts a time struct (e.g. time.Time) to pgtype.Timestamptz.
func Timestamptz(v time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: v, Status: pgtype.Present}
}
