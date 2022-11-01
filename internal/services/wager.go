package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/wager-api/internal/entities"
	"github.com/wager-api/internal/models"
	"github.com/wager-api/internal/repositories"
	"github.com/wager-api/libs/database"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"go.uber.org/multierr"
)

type WagerService struct {
	DB        database.Ext
	WagerRepo interface {
		Create(ctx context.Context, db database.Ext, wager *entities.Wager) error
		Update(ctx context.Context, db database.Ext, wager *entities.Wager) (pgconn.CommandTag, error)
		Get(ctx context.Context, db database.Ext, wagerID pgtype.Int4, queryEnhancers ...repositories.QueryEnhancer) (*entities.Wager, error)
		List(ctx context.Context, db database.Ext, lastID pgtype.Int4, limit uint32) ([]*entities.Wager, error)
	}
	PurchaseRepo interface {
		Create(ctx context.Context, db database.Ext, purchase *entities.Purchase) error
	}
}

func validatePlaceWagerReq(req *models.PlaceWagerRequest) error {
	if req.TotalWagerValue <= 0 {
		return fmt.Errorf("the total_wager_value must be a positive integer above 0")
	}
	if req.Odds <= 0 {
		return fmt.Errorf("the odds must be a positive integer above 0")
	}
	if req.SellingPercentage < 1 || req.SellingPercentage > 100 {
		return fmt.Errorf("the selling_percentage must be specified as an integer between 1 and 100")
	}
	tempSellingPriceNumber := req.SellingPrice * 100
	if req.SellingPrice <= 0 || tempSellingPriceNumber-float32(int(tempSellingPriceNumber)) > 0 {
		return fmt.Errorf("the selling_price must be a positive decimal value to two decimal places")
	}
	if req.SellingPrice <= req.TotalWagerValue*float32(req.SellingPercentage)/100 {
		return fmt.Errorf("selling_price must be greater than total_wager_value * (selling_percentage / 100)")
	}

	return nil
}
func (s *WagerService) PlaceWager(resp http.ResponseWriter, req *http.Request) {
	placeWagerRequest := &models.PlaceWagerRequest{}
	err := json.NewDecoder(req.Body).Decode(&placeWagerRequest)
	defer req.Body.Close()
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(resp).Encode(map[string]string{
			"error": "unable to parse request",
		})
		return
	}
	if err := validatePlaceWagerReq(placeWagerRequest); err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(resp).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	ctx := req.Context()

	wager := &entities.Wager{}
	now := time.Now()
	database.AllNullEntity(wager)
	if err = multierr.Combine(
		wager.TotalWagerValue.Set(placeWagerRequest.TotalWagerValue),
		wager.Odds.Set(placeWagerRequest.Odds),
		wager.SellingPercentage.Set(placeWagerRequest.SellingPercentage),
		wager.SellingPrice.Set(placeWagerRequest.SellingPrice),
		wager.CurrentSellingPrice.Set(placeWagerRequest.SellingPrice),
		wager.PlaceAt.Set(now),
		wager.CreatedAt.Set(now),
		wager.UpdatedAt.Set(now),
	); err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(resp).Encode(map[string]string{
			"error": "unable to generate value for wager",
		})
		return
	}
	if err := s.WagerRepo.Create(ctx, s.DB, wager); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(resp).Encode(map[string]string{
			"error": "unable to create wager",
		})
		return
	}
	resp.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(resp).Encode(convertWagerPg2placeWagerResponse(wager))
}
func convertWagerPg2placeWagerResponse(wager *entities.Wager) *models.PlaceWagerResponse {
	var placedAt *time.Time
	if wager.PlaceAt.Status == pgtype.Present {
		placedAt = &wager.PlaceAt.Time
	}
	return &models.PlaceWagerResponse{
		ID:                  int(wager.WagerID.Int),
		TotalWagerValue:     wager.TotalWagerValue.Float,
		Odds:                int(wager.Odds.Int),
		SellingPercentage:   int(wager.SellingPercentage.Int),
		SellingPrice:        wager.SellingPrice.Float,
		CurrentSellingPrice: wager.CurrentSellingPrice.Float,
		PercentageSold:      wager.PercentageSold.Float,
		AmountSold:          wager.AmountSold.Float,
		PlacedAt:            placedAt,
	}
}

func convertWagerPg2Domain(wager *entities.Wager) *models.Wager {
	var placedAt *time.Time
	if wager.PlaceAt.Status == pgtype.Present {
		placedAt = &wager.PlaceAt.Time
	}
	return &models.Wager{
		ID:                  int(wager.WagerID.Int),
		TotalWagerValue:     wager.TotalWagerValue.Float,
		Odds:                int(wager.Odds.Int),
		SellingPercentage:   int(wager.SellingPercentage.Int),
		SellingPrice:        wager.SellingPrice.Float,
		CurrentSellingPrice: wager.CurrentSellingPrice.Float,
		PercentageSold:      wager.PercentageSold.Float,
		AmountSold:          wager.AmountSold.Float,
		PlacedAt:            placedAt,
	}
}

func validateBuyWagerReq(req *models.BuyWagerRequest) error {
	if req.BuyingPrice <= 0 {
		return fmt.Errorf("the buying_price must be a positive decimal")
	}
	return nil
}
func (s *WagerService) BuyWager(resp http.ResponseWriter, req *http.Request) {
	buyWagerRequest := &models.BuyWagerRequest{}
	err := json.NewDecoder(req.Body).Decode(&buyWagerRequest)
	defer req.Body.Close()
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(resp).Encode(map[string]string{
			"error": "unable to parse request",
		})
		return
	}
	if err := validateBuyWagerReq(buyWagerRequest); err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(resp).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	ctx := req.Context()
	var wagerID int
	if wagerIDRaw, ok := ctx.Value("wager_id").(int); ok {
		wagerID = wagerIDRaw
	} else {
		resp.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(resp).Encode(map[string]string{
			"error": "wager_id wrong format",
		})
		return
	}
	purchaseRecord := &entities.Purchase{}
	database.AllNullEntity(purchaseRecord)
	//TODO: validate req
	if err := database.ExecInTx(ctx, s.DB, func(ctx context.Context, tx pgx.Tx) error {
		wager, err := s.WagerRepo.Get(ctx, tx, database.Int4(int32((wagerID))), repositories.WithUpdateLock())
		if err != nil {
			if err == pgx.ErrNoRows {
				return fmt.Errorf("unable to get wager information: not found")
			}
			return fmt.Errorf("unable to get wager information")
		}
		if buyWagerRequest.BuyingPrice > wager.CurrentSellingPrice.Float {
			return fmt.Errorf("unable to execute: buying_price must be lesser or equal to current_selling_price")
		}
		now := time.Now()

		if err = multierr.Combine(
			purchaseRecord.WagerID.Set(wagerID),
			purchaseRecord.BuyingPrice.Set(buyWagerRequest.BuyingPrice),
			purchaseRecord.BoughtAt.Set(now),
			purchaseRecord.CreatedAt.Set(now),
			purchaseRecord.UpdatedAt.Set(now)); err != nil {
			return fmt.Errorf("unable to generate new purchase record")
		}
		err = s.PurchaseRepo.Create(ctx, tx, purchaseRecord)
		if err != nil {
			return fmt.Errorf("unable to create new purchase record")
		}
		if err = multierr.Combine(
			wager.CurrentSellingPrice.Set(buyWagerRequest.BuyingPrice),
			wager.AmountSold.Set(wager.AmountSold.Float+buyWagerRequest.BuyingPrice),
			wager.PercentageSold.Set(roundFloat((wager.AmountSold.Float/wager.SellingPrice.Float)*100)),

			wager.UpdatedAt.Set(now)); err != nil {
			return fmt.Errorf("unable to generate wager record")
		}

		cmdTag, err := s.WagerRepo.Update(ctx, tx, wager)
		if err != nil {

			return fmt.Errorf("unable to update wager record")
		}
		if cmdTag.RowsAffected() != 1 {
			return fmt.Errorf("unable to update wager record: no row affected")
		}
		return nil
	}); err != nil {
		if err.Error() == "unable to get wager information: not found" {
			resp.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(resp).Encode(map[string]string{
				"error": fmt.Sprintf("unable to buy wager: %s", err),
			})
			return
		}
		resp.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(resp).Encode(map[string]string{
			"error": fmt.Sprintf("unable to buy wager: %s", err.Error()),
		})
		return
	}
	resp.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(resp).Encode(convert2BuyWagerResponse(purchaseRecord))
}

func convert2BuyWagerResponse(purchase *entities.Purchase) *models.BuyWagerResponse {
	return &models.BuyWagerResponse{
		PurchaseID:  int(purchase.PurchaseID.Int),
		WagerID:     int(purchase.WagerID.Int),
		BuyingPrice: purchase.BuyingPrice.Float,
		BoughtAt:    &purchase.BoughtAt.Time,
	}
}

// roundFloat ensure round to two decimal places
func roundFloat(number float32) float32 {
	return float32(math.Round((float64(number) * 100)) / 100)
}

func validateListWagerParam(req *http.Request) error {

	ctx := req.Context()
	var page, limit int

	if pageRaw, ok := ctx.Value("page").(int); ok {
		page = pageRaw
	} else {
		return fmt.Errorf("page wrong format")
	}
	if limitRaw, ok := ctx.Value("limit").(int); ok {
		limit = limitRaw
	} else {
		return fmt.Errorf("limit wrong format")
	}
	if page <= 0 || limit <= 0 {
		return fmt.Errorf("`page` must be positive number and `limit` should be greater than 0")
	}
	return nil
}

func (s *WagerService) ListWager(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	if err := validateListWagerParam(req); err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(resp).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	page, limit := ctx.Value("page").(int), ctx.Value("limit").(int)
	var lastID pgtype.Int4
	_ = lastID.Set(nil)
	if page > 1 {
		_ = lastID.Set((page - 1) * limit)
	}
	wagers, err := s.WagerRepo.List(ctx, s.DB, lastID, uint32(limit))
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(resp).Encode(map[string]string{
			"error": "unable to list wager",
		})
		return
	}
	wagermodels := make([]*models.Wager, 0, len(wagers))
	for _, wager := range wagers {
		wagermodels = append(wagermodels, convertWagerPg2Domain(wager))
	}
	resp.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(resp).Encode(wagermodels)
}
