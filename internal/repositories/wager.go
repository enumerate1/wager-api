package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/wager-api/internal/entities"
	"github.com/wager-api/libs/database"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
)

type WagerRepo struct{}

func (r *WagerRepo) Create(ctx context.Context, db database.Ext, wager *entities.Wager) error {
	command := `INSERT INTO %s (%s) VALUES (%s) RETURNING wager_id`
	fieldNames := database.GetFieldNamesExcepts(wager, []string{"wager_id"})
	fmt.Println("===internal", fieldNames, len(fieldNames))
	placeHolders := database.GeneratePlaceholders(len(fieldNames))
	ultimateCmd := fmt.Sprintf(command, wager.TableName(), strings.Join(fieldNames, ","), placeHolders)
	args := database.GetScanFields(wager, fieldNames)
	fmt.Println("args", args)
	if err := db.QueryRow(ctx, ultimateCmd, args...).Scan(&wager.WagerID); err != nil {
		return err
	}
	return nil
}

func (r *WagerRepo) Update(ctx context.Context, db database.Ext, wager *entities.Wager) (pgconn.CommandTag, error) {
	query := fmt.Sprintf(
		`
		   UPDATE %s
		   SET current_selling_price = $1, percentage_sold = $2, amount_sold = $3, updated_at = now()
		   WHERE
		     wager_id = $4 AND
		     deleted_at IS NULL
	       `,
		wager.TableName(),
	)
	cmdTag, err := db.Exec(ctx, query, wager.CurrentSellingPrice, wager.PercentageSold, wager.AmountSold, wager.WagerID)
	if err != nil {
		return cmdTag, fmt.Errorf("db.Exec: %w", err)
	}

	return cmdTag, nil
}

func (r *WagerRepo) Get(ctx context.Context, db database.Ext, wagerID pgtype.Int4, queryEnhancers ...QueryEnhancer) (*entities.Wager, error) {
	getWagerCmd := `SELECT %s FROM %s WHERE wager_id = $1 AND deleted_at IS NULL`
	wagerEnt := &entities.Wager{}
	fields, values := wagerEnt.FieldMap()
	for _, e := range queryEnhancers {
		e(&getWagerCmd)
	}

	err := db.QueryRow(ctx, fmt.Sprintf(getWagerCmd, strings.Join(fields, ", "), wagerEnt.TableName()), &wagerID).Scan(values...)
	if err != nil {
		return nil, err
	}

	return wagerEnt, nil
}

func (r *WagerRepo) List(ctx context.Context, db database.Ext, lastID pgtype.Int4, limit uint32) ([]*entities.Wager, error) {
	b := &entities.Wager{}
	fieldName, _ := b.FieldMap()
	query := fmt.Sprintf("SELECT %s FROM %s WHERE ($1::INT IS NULL OR wager_id>$1) ORDER BY wager_id LIMIT $2 ", strings.Join(fieldName, ", "), b.TableName())
	wagers := entities.Wagers{}
	if err := database.Select(ctx, db, query, lastID, limit).ScanAll(&wagers); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return wagers, nil
}
