package entities

import (
	"github.com/wager-api/libs/database"

	"github.com/jackc/pgtype"
)

type Wager struct {
	WagerID             pgtype.Int4
	TotalWagerValue     pgtype.Float4
	Odds                pgtype.Int4
	SellingPercentage   pgtype.Int4
	SellingPrice        pgtype.Float4
	CurrentSellingPrice pgtype.Float4
	PercentageSold      pgtype.Float4
	AmountSold          pgtype.Float4
	PlaceAt             pgtype.Timestamptz
	CreatedAt           pgtype.Timestamptz
	UpdatedAt           pgtype.Timestamptz
	DeletedAt           pgtype.Timestamptz
}

func (e *Wager) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		"wager_id",
		"total_wager_value",
		"odds",
		"selling_percentage",
		"selling_price",
		"current_selling_price",
		"percentage_sold",
		"amount_sold",
		"place_at",
		"created_at",
		"updated_at",
		"deleted_at",
	}
	values = []interface{}{
		&e.WagerID,
		&e.TotalWagerValue,
		&e.Odds,
		&e.SellingPercentage,
		&e.SellingPrice,
		&e.CurrentSellingPrice,
		&e.PercentageSold,
		&e.AmountSold,
		&e.PlaceAt,
		&e.CreatedAt,
		&e.UpdatedAt,
		&e.DeletedAt,
	}
	return
}
func (e *Wager) TableName() string {
	return "wager"
}

type Wagers []*Wager

func (es *Wagers) Add() database.Entity {
	e := &Wager{}
	*es = append(*es, e)
	return e
}
