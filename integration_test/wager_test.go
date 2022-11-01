package integrationtest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/wager-api/internal/entities"
	"github.com/wager-api/internal/models"
	"github.com/wager-api/internal/repositories"
	"github.com/wager-api/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"github.com/wager-api/libs/database"
	"github.com/wager-api/libs/mux"

	"go.uber.org/zap"
)

var (
	DB     database.Ext
	chiMux *chi.Mux
)

func Test_PlaceWager(t *testing.T) {
	var tests = []struct {
		description    string
		url            string
		jsonRequest    []byte
		expectedStatus int
		expecteValue   models.PlaceWagerResponse
	}{
		{
			"happy case",
			"/wagers",
			[]byte(`{"total_wager_value": 50, "odds": 30,"selling_percentage": 30,"selling_price": 50}`),
			201,
			models.PlaceWagerResponse{
				TotalWagerValue:     50,
				Odds:                30,
				SellingPercentage:   30,
				SellingPrice:        50,
				CurrentSellingPrice: 50,
			},
		},
	}
	for _, tt := range tests {
		t.Run(strings.Join([]string{tt.url, tt.description}, "_"), func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.url, bytes.NewBuffer([]byte(tt.jsonRequest)))
			rec := httptest.NewRecorder()
			chiMux.ServeHTTP(rec, req)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus == http.StatusCreated {
				placeWagerResponse := models.PlaceWagerResponse{}
				err := json.NewDecoder(rec.Body).Decode(&placeWagerResponse)
				defer req.Body.Close()
				assert.Equal(t, tt.expectedStatus, rec.Code, "status code must be 201")
				assert.NoError(t, err)
				assert.NotEqual(t, 0, placeWagerResponse.ID, "wager_id must be generated")
				assert.NotEqual(t, nil, placeWagerResponse.PlacedAt, "placed_at must be have value")
				assert.Equal(t, placeWagerResponse.SellingPrice, placeWagerResponse.CurrentSellingPrice, "selling_price and current_selling_price must be equal in first initiation")
			}

		})
	}

}

// Test_BuyWager_HappyCase
// Step 1: init Wager by call PlaceWager
// Step 2: call BuyWager
// Step 3: check BuyWager and the wager's info in DB
func Test_BuyWager_HappyCase(t *testing.T) {
	// Step 1: init Wager by call PlaceWager
	placeWagerDataReq := []byte(`{"total_wager_value": 50, "odds": 30,"selling_percentage": 30,"selling_price": 50}`)
	placeWagerReq := httptest.NewRequest(http.MethodPost, "/wagers", bytes.NewBuffer([]byte(placeWagerDataReq)))
	rec := httptest.NewRecorder()
	chiMux.ServeHTTP(rec, placeWagerReq)
	assert.Equal(t, 201, rec.Code, "status code must be 201")
	placeWagerResponse := models.PlaceWagerResponse{}
	err := json.NewDecoder(rec.Body).Decode(&placeWagerResponse)
	defer placeWagerReq.Body.Close()
	assert.NoError(t, err)

	// Step 2: call BuyWager
	buyWagerDataReq := models.BuyWagerRequest{
		BuyingPrice: 40,
	}
	buyWagerDataReqByte, err := json.Marshal(buyWagerDataReq)
	assert.NoError(t, err)
	buyWagerReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/buy/%d", placeWagerResponse.ID), bytes.NewBuffer([]byte(buyWagerDataReqByte)))
	rec = httptest.NewRecorder()
	chiMux.ServeHTTP(rec, buyWagerReq)

	// Step 3: check BuyWager and the wager's info in DB
	assert.Equal(t, 201, rec.Code, "status code must be 201")
	buyWagerResponse := models.BuyWagerResponse{}
	err = json.NewDecoder(rec.Body).Decode(&buyWagerResponse)
	assert.NoError(t, err)
	defer placeWagerReq.Body.Close()
	assert.Equal(t, placeWagerResponse.ID, buyWagerResponse.WagerID, "2 wager ID must be the same")
	// Step 3 (plus) check the wager's info in DB
	wagerEnt := &entities.Wager{}
	fieldNames, fields := wagerEnt.FieldMap()
	ctx := context.Background()
	cmd := fmt.Sprintf(`SELECT %s FROM wager WHERE wager_id = $1`, strings.Join(fieldNames, ","))
	err = DB.QueryRow(ctx, cmd, placeWagerResponse.ID).Scan(fields...)
	assert.NoError(t, err)

	// first time so AmountSold = CurrentSellingPrice
	assert.Equal(t, wagerEnt.CurrentSellingPrice.Float, wagerEnt.AmountSold.Float)
	assert.Equal(t, roundFloat((wagerEnt.AmountSold.Float/wagerEnt.SellingPrice.Float)*100), wagerEnt.PercentageSold.Float)
}

// roundFloat ensure round to two decimal places
func roundFloat(number float32) float32 {
	return float32(math.Round((float64(number) * 100)) / 100)
}

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Printf("Could not start new pool: %s", err)
	}
	resource, err := pool.Run("postgres", "14-alpine", []string{"POSTGRES_PASSWORD=password", "POSTGRES_DB=postgres"})

	if err != nil {
		log.Printf("Could not start resource: %s", err)
	}
	if err = pool.Retry(func() error {
		pool, err := pgxpool.Connect(context.Background(), fmt.Sprintf("postgresql://postgres:password@localhost:%s/%s?sslmode=disable", resource.GetPort("5432/tcp"), "postgres"))
		if err != nil {
			log.Printf("Could not connect resource: %s", err)
			return err
		}
		migrate, err := migrate.New(
			"file://./../postgres", // depend on your migrations
			fmt.Sprintf("postgres://postgres:password@localhost:%s/%s?sslmode=disable", resource.GetPort("5432/tcp"), "postgres"),
		)
		if err := migrate.Up(); err != nil {
			// log.Println(err)
			// log.Println(os.Getwd())

			return err
		}
		wagerService := &services.WagerService{
			DB:           pool,
			WagerRepo:    &repositories.WagerRepo{},
			PurchaseRepo: &repositories.PurchaseRepo{},
		}
		DB = pool

		chiMux = mux.InitWithLogger((zap.NewNop()))
		services.NewWagerHandler(chiMux, wagerService)
		if err != nil {
			log.Print("Could not migrate", err)
			return err
		}

		return nil
	}); err != nil {
		log.Printf("Could not connect to docker: %s", err)
	}
	if err != nil {
		log.Printf("Could not connect resource: %s", err)
	}
	code := m.Run()
	// When you're done, kill and remove the container
	if err = pool.Purge(resource); err != nil {
		log.Printf("Could not purge resource: %s", err)
	}
	os.Exit(code)
}
