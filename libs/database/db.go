package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type Row interface {
	pgx.Row
}

type Rows interface {
	pgx.Rows
}

// queryer is an interface for Query
type queryer interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
}

// execer is an interfacr for Exec
type execer interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

type QueryExecer interface {
	queryer
	execer
}

type TxStarter interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type TxController interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// Ext is an union interface which can bind query and exec
type Ext interface {
	QueryExecer
	TxStarter
}
type TxHandler = func(ctx context.Context, tx pgx.Tx) error

func ExecInTx(ctx context.Context, db Ext, txHandler TxHandler) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("db.Begin: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}
		err = fmt.Errorf("tx.Commit: %w", tx.Commit(ctx))
	}()
	err = txHandler(ctx, tx)
	return err
}

type Tx interface {
	pgx.Tx
}
