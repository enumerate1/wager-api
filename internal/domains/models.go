package domains

import "time"

type Wager struct {
	ID                  int        `json:"id,omitempty"`
	TotalWagerValue     float32    `json:"total_wager_value,omitempty" validate:"gt=0"`
	Odds                int        `json:"odds,omitempty" validate:"gt=0"`
	SellingPercentage   int        `json:"selling_percentage,omitempty" validate:"gte=1,lte=100"`
	SellingPrice        float32    `json:"selling_price,omitempty" validate:"gte=1,lte=100"`
	CurrentSellingPrice float32    `json:"current_selling_price,omitempty"`
	PercentageSold      float32    `json:"percentage_sold,omitempty"`
	AmountSold          float32    `json:"amount_sold,omitempty"`
	PlacedAt            *time.Time `json:"placed_at,omitempty"`
}

type PlaceWagerRequest struct {
	TotalWagerValue   float32 `json:"total_wager_value" validate:"gt=0"`
	Odds              int     `json:"odds" validate:"gt=0"`
	SellingPercentage int     `json:"selling_percentage" validate:"gte=1,lte=100"`
	SellingPrice      float32 `json:"selling_price" validate:"gte=1,lte=100"`
}

type PlaceWagerResponse struct {
	ID                  int        `json:"id"`
	TotalWagerValue     float32    `json:"total_wager_value" validate:"gt=0"`
	Odds                int        `json:"odds" validate:"gt=0"`
	SellingPercentage   int        `json:"selling_percentage" validate:"gte=1,lte=100"`
	SellingPrice        float32    `json:"selling_price" validate:"gte=1,lte=100"`
	CurrentSellingPrice float32    `json:"current_selling_price"`
	PercentageSold      float32    `json:"percentage_sold"`
	AmountSold          float32    `json:"amount_sold"`
	PlacedAt            *time.Time `json:"placed_at"`
}

type BuyWagerRequest struct {
	BuyingPrice float32 `json:"buying_price,omitempty"`
}
type BuyWagerResponse struct {
	PurchaseID  int        `json:"purchase_id"`
	WagerID     int        `json:"wager_id"`
	BuyingPrice float32    `json:"buying_price"`
	BoughtAt    *time.Time `json:"bought_at"`
}
