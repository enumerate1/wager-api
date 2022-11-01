package entities

import (
	"github.com/wager-api/libs/database"

	"github.com/jackc/pgtype"
)

type Purchase struct {
	PurchaseID  pgtype.Int4
	WagerID     pgtype.Int4
	BuyingPrice pgtype.Float4
	BoughtAt    pgtype.Timestamptz
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
	DeletedAt   pgtype.Timestamptz
}

func (e *Purchase) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		"purchase_id",
		"wager_id",
		"buying_price",
		"bought_at",
		"created_at",
		"updated_at",
		"deleted_at",
	}
	values = []interface{}{
		&e.PurchaseID,
		&e.WagerID,
		&e.BuyingPrice,
		&e.BoughtAt,
		&e.CreatedAt,
		&e.UpdatedAt,
		&e.UpdatedAt,
		&e.DeletedAt,
	}
	return
}
func (e *Purchase) TableName() string {
	return "purchase"
}

type Purchases []*Purchase

func (es *Purchases) Add() database.Entity {
	e := &Purchase{}
	*es = append(*es, e)
	return e
}
