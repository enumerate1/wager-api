package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/wager-api/internal/entities"
	"github.com/wager-api/libs/database"
)

type PurchaseRepo struct{}

func (r *PurchaseRepo) Create(ctx context.Context, db database.Ext, purchase *entities.Purchase) error {
	command := `INSERT INTO %s (%s) VALUES (%s) RETURNING purchase_id`
	// fieldNames, _ := purchase.FieldMap()
	// fields := []string{"wager_id",
	// 	"buying_price",
	// 	"bought_at",
	// 	"created_at",
	// 	"updated_at",
	// 	"deleted_at"}
	fieldNames := database.GetFieldNamesExcepts(purchase, []string{"purchase_id"})
	placeHolders := database.GeneratePlaceholders(len(fieldNames))
	ultimateCmd := fmt.Sprintf(command, purchase.TableName(), strings.Join(fieldNames, ","), placeHolders)
	args := database.GetScanFields(purchase, fieldNames)
	if err := db.QueryRow(ctx, ultimateCmd, args...).Scan(&purchase.PurchaseID); err != nil {
		return err
	}
	return nil
}
