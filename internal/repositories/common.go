package repositories

import (
	"errors"

	"github.com/jackc/pgx/v4"
)

var (
	ErrNoRows          = pgx.ErrNoRows
	ErrUniqueViolation = errors.New("unique violation.")
)

type QueryEnhancer func(query *string)

func WithUpdateLock() QueryEnhancer {
	return func(query *string) {
		*query += " FOR UPDATE"
	}
}
