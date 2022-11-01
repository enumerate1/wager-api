package services

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wager-api/internal/entities"
	"github.com/wager-api/internal/models"
	"github.com/wager-api/libs/database"
	mock_database "github.com/wager-api/mocks/libs/database"
	mock_repositories "github.com/wager-api/mocks/repositories"

	"github.com/jackc/pgtype"
)

type TestCase struct {
	ctx            context.Context
	name           string
	jsonReq        []byte
	url            string
	expectedResp   interface{}
	expectedStatus int

	setup func(ctx context.Context)
}

func Test_validatePlaceWagerReq(t *testing.T) {
	t.Parallel()
	type testcase struct {
		name          string
		placeWagerReq *models.PlaceWagerRequest
		expectedErr   error
	}
	tests := []testcase{
		{
			name:        "total_wager_value = 0",
			expectedErr: fmt.Errorf("the total_wager_value must be a positive integer above 0"),
			placeWagerReq: &models.PlaceWagerRequest{
				TotalWagerValue: 0,
			},
		},
		{
			name:        "odds = 0",
			expectedErr: fmt.Errorf("the odds must be a positive integer above 0"),
			placeWagerReq: &models.PlaceWagerRequest{
				TotalWagerValue: 10,
				Odds:            0,
			},
		},
		{
			name:        "selling_percentage out of range [1, 100]",
			expectedErr: fmt.Errorf("the selling_percentage must be specified as an integer between 1 and 100"),
			placeWagerReq: &models.PlaceWagerRequest{
				TotalWagerValue:   10,
				Odds:              20,
				SellingPercentage: 110,
			},
		},
		{
			name:        "selling_price have more than 2 decimal places out of range [1, 100]",
			expectedErr: fmt.Errorf("the selling_price must be a positive decimal value to two decimal places"),
			placeWagerReq: &models.PlaceWagerRequest{
				TotalWagerValue:   10,
				Odds:              20,
				SellingPercentage: 50,
				SellingPrice:      10.112,
			},
		},
		{
			name:        "selling_price lesser than total_wager_value * (selling_percentage / 100)",
			expectedErr: fmt.Errorf("selling_price must be greater than total_wager_value * (selling_percentage / 100)"),
			placeWagerReq: &models.PlaceWagerRequest{
				TotalWagerValue:   10,
				Odds:              20,
				SellingPercentage: 50,
				SellingPrice:      4,
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := validatePlaceWagerReq(tc.placeWagerReq)
			assert.Equal(t, tc.expectedErr, err)

		})
	}
}

func Test_validateBuyWagerReq(t *testing.T) {
	// t.Parallel()
	type testcase struct {
		name        string
		buyWagerReq *models.BuyWagerRequest
		expectedErr error
	}
	tests := []testcase{
		{
			name:        "buying_price less or equal to 0",
			expectedErr: fmt.Errorf("the buying_price must be a positive decimal"),
			buyWagerReq: &models.BuyWagerRequest{
				BuyingPrice: 0,
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := validateBuyWagerReq(tc.buyWagerReq)
			fmt.Println("===err", err)
			assert.Equal(t, tc.expectedErr, err)

		})
	}
}
func Test_validateListWagerParams(t *testing.T) {
	type testcase struct {
		name        string
		req         *http.Request
		expectedErr error
	}
	// reqWithZeroValue
	ctx := context.Background()
	reqWithZeroValue := &http.Request{}
	ctx = context.WithValue(ctx, "limit", 0)
	ctx = context.WithValue(ctx, "page", 0)
	reqWithZeroValue = reqWithZeroValue.WithContext(ctx)
	//reqWithWrongFormatParam
	ctxreqWithWrongFormatParam := context.Background()
	reqWithWrongFormatParam := &http.Request{}
	ctxreqWithWrongFormatParam = context.WithValue(ctxreqWithWrongFormatParam, "limit", "")
	ctxreqWithWrongFormatParam = context.WithValue(ctxreqWithWrongFormatParam, "page", "")
	reqWithWrongFormatParam = reqWithZeroValue.WithContext(ctxreqWithWrongFormatParam)

	tests := []testcase{
		{
			name:        "limit = 0, page = 0",
			expectedErr: fmt.Errorf("`page` must be positive number and `limit` should be greater than 0"),
			req:         reqWithZeroValue,
		},
		{
			name:        "req with wrong format param",
			expectedErr: fmt.Errorf("page wrong format"),
			req:         reqWithWrongFormatParam,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := validateListWagerParam(tc.req)
			assert.Equal(t, tc.expectedErr, err)

		})
	}
}

func Test_PlaceWager(t *testing.T) {
	t.Parallel()

	db := &mock_database.Ext{}
	wagerRepo := &mock_repositories.MockWagerRepo{}
	purchaseRepo := &mock_repositories.MockPurchaseRepo{}
	mockErr := fmt.Errorf("mock-error")
	testcases := []TestCase{
		{
			// validation request
			name:           "bad request (violate input condition)",
			expectedResp:   []byte(`{"error":"the total_wager_value must be a positive integer above 0"}`),
			url:            "/wagers",
			jsonReq:        []byte(`{"total_wager_value": 0, "odds": 30,"selling_percentage": 30,"selling_price": 50}`),
			expectedStatus: http.StatusBadRequest,
			setup: func(ctx context.Context) {
			},
		},
		{
			name:           "err when create wager",
			expectedResp:   []byte(`{"error":"unable to create wager"}`),
			url:            "/wagers",
			jsonReq:        []byte(`{"total_wager_value": 20, "odds": 30,"selling_percentage": 30,"selling_price": 50}`),
			expectedStatus: http.StatusInternalServerError,
			setup: func(ctx context.Context) {
				wagerRepo.On("Create", ctx, db, mock.Anything, mock.Anything).Once().Return(mockErr)
			},
		},
	}

	wagerService := &WagerService{
		DB:           db,
		WagerRepo:    wagerRepo,
		PurchaseRepo: purchaseRepo,
	}
	mockWagerHandler := WagerHandler{WagerService: wagerService}
	for _, tc := range testcases {
		t.Run(tc.url, func(t *testing.T) {
			tc.setup(context.Background())
			req := httptest.NewRequest(http.MethodPost, tc.url, bytes.NewBuffer([]byte(tc.jsonReq)))
			rec := httptest.NewRecorder()
			http.HandlerFunc(mockWagerHandler.WagerService.PlaceWager).ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedStatus, rec.Code)
			data, err := ioutil.ReadAll(rec.Body)
			assert.NoError(t, err)
			// :len(data)-1 remove the `\n`
			assert.Equal(t, tc.expectedResp, data[:len(data)-1])

		})
	}
}
func Test_BuyWager(t *testing.T) {
	t.Parallel()
	db := &mock_database.Ext{}
	tx := &mock_database.Tx{}
	wagerRepo := &mock_repositories.MockWagerRepo{}
	purchaseRepo := &mock_repositories.MockPurchaseRepo{}
	ctx := context.Background()
	wagerID := 1
	ctx = context.WithValue(ctx, "wager_id", wagerID)
	mockErr := fmt.Errorf("mock-error")
	testcases := []TestCase{
		{
			ctx:            ctx,
			name:           "error when get the wager information",
			expectedResp:   []byte(`{"error":"unable to buy wager: unable to get wager information"}`),
			url:            "/buy/1",
			jsonReq:        []byte(`{"buying_price": 20}`),
			expectedStatus: http.StatusInternalServerError,
			setup: func(ctx context.Context) {
				db.On("Begin", ctx).Return(tx, nil)
				tx.On("Rollback", mock.Anything).Return(nil)
				wagerRepo.On("Get", ctx, tx, database.Int4(int32(wagerID))).Once().Return(nil, mockErr)
			},
		},
		{
			// validation request
			name:           "bad request (violate input condition)",
			expectedResp:   []byte(`{"error":"the buying_price must be a positive decimal"}`),
			url:            "/wagers",
			jsonReq:        []byte(`{"buying_price": 0}`),
			expectedStatus: http.StatusBadRequest,
			setup: func(ctx context.Context) {
			},
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.url, func(t *testing.T) {
			tc.setup(ctx)
			wagerService := &WagerService{
				DB:           db,
				WagerRepo:    wagerRepo,
				PurchaseRepo: purchaseRepo,
			}
			mockWagerHandler := WagerHandler{WagerService: wagerService}
			req := httptest.NewRequest(http.MethodPost, tc.url, bytes.NewBuffer([]byte(tc.jsonReq)))
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()
			http.HandlerFunc(mockWagerHandler.WagerService.BuyWager).ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedStatus, rec.Code)
			data, err := ioutil.ReadAll(rec.Body)
			assert.NoError(t, err)
			// :len(data)-1 remove the `\n`
			assert.Equal(t, tc.expectedResp, data[:len(data)-1])

		})
	}
}

func Test_ListWager(t *testing.T) {
	t.Parallel()

	db := &mock_database.Ext{}
	wagerRepo := &mock_repositories.MockWagerRepo{}
	purchaseRepo := &mock_repositories.MockPurchaseRepo{}
	ctx := context.Background()
	page := 1
	limit := 10
	nilWagerID := pgtype.Int4{}
	_ = nilWagerID.Set(nil)
	ctx = context.WithValue(ctx, "page", page)
	ctx = context.WithValue(ctx, "limit", limit)
	wagers := []*entities.Wager{
		{
			WagerID:         database.Int4(1),
			TotalWagerValue: database.Float4(100),
		},
		{
			WagerID:         database.Int4(2),
			TotalWagerValue: database.Float4(100),
		},
	}
	mockErr := fmt.Errorf("mock-error")
	testcases := []TestCase{
		{
			ctx:          ctx,
			name:         "happy case",
			expectedResp: []byte(`[{"id":1,"total_wager_value":100,"odds":0,"selling_percentage":0,"selling_price":0,"current_selling_price":0,"percentage_sold":0,"amount_sold":0,"placed_at":null},{"id":2,"total_wager_value":100,"odds":0,"selling_percentage":0,"selling_price":0,"current_selling_price":0,"percentage_sold":0,"amount_sold":0,"placed_at":null}]`),
			// work in both cases /wagers?page=:4&limit=:4 and /wagers?page=4&limit=4
			url:            "/wagers?page=:4&limit=:4",
			expectedStatus: http.StatusOK,
			setup: func(ctx context.Context) {
				wagerRepo.On("List", ctx, db, nilWagerID, uint32(limit)).Once().Return(wagers, nil)
			},
		},
		{
			ctx:          ctx,
			name:         "error when list the wagers",
			expectedResp: []byte(`{"error":"unable to list wager"}`),
			// work in both cases /wagers?page=:4&limit=:4 and /wagers?page=4&limit=4
			url:            "/wagers?page=:4&limit=:4",
			expectedStatus: http.StatusInternalServerError,
			setup: func(ctx context.Context) {
				wagerRepo.On("List", ctx, db, nilWagerID, uint32(limit)).Once().Return(nil, mockErr)
			},
		},
	}
	wagerService := &WagerService{
		DB:           db,
		WagerRepo:    wagerRepo,
		PurchaseRepo: purchaseRepo,
	}
	mockWagerHandler := WagerHandler{WagerService: wagerService}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.url, func(t *testing.T) {
			tc.setup(ctx)

			req := httptest.NewRequest(http.MethodGet, tc.url, bytes.NewBuffer([]byte(tc.jsonReq)))
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()
			http.HandlerFunc(mockWagerHandler.WagerService.ListWager).ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedStatus, rec.Code)
			data, err := ioutil.ReadAll(rec.Body)
			assert.NoError(t, err)
			// :len(data)-1 remove the `\n`
			assert.Equal(t, tc.expectedResp, data[:len(data)-1])

		})
	}
}
