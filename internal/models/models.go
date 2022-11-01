package models

import "time"

type Wager struct {
	ID                  int        `json:"id,omitempty"`
	TotalWagerValue     float32    `json:"total_wager_value"`
	Odds                int        `json:"odds"`
	SellingPercentage   int        `json:"selling_percentage"`
	SellingPrice        float32    `json:"selling_price"`
	CurrentSellingPrice float32    `json:"current_selling_price"`
	PercentageSold      float32    `json:"percentage_sold"`
	AmountSold          float32    `json:"amount_sold"`
	PlacedAt            *time.Time `json:"placed_at"`
}

type PlaceWagerRequest struct {
	TotalWagerValue   float32 `json:"total_wager_value"`
	Odds              int     `json:"odds"`
	SellingPercentage int     `json:"selling_percentage"`
	SellingPrice      float32 `json:"selling_price"`
}

type PlaceWagerResponse struct {
	ID                  int        `json:"id"`
	TotalWagerValue     float32    `json:"total_wager_value"`
	Odds                int        `json:"odds"`
	SellingPercentage   int        `json:"selling_percentage"`
	SellingPrice        float32    `json:"selling_price"`
	CurrentSellingPrice float32    `json:"current_selling_price"`
	PercentageSold      float32    `json:"percentage_sold"`
	AmountSold          float32    `json:"amount_sold"`
	PlacedAt            *time.Time `json:"placed_at"`
}

type BuyWagerRequest struct {
	BuyingPrice float32 `json:"buying_price"`
}
type BuyWagerResponse struct {
	PurchaseID  int        `json:"purchase_id"`
	WagerID     int        `json:"wager_id"`
	BuyingPrice float32    `json:"buying_price"`
	BoughtAt    *time.Time `json:"bought_at"`
}
