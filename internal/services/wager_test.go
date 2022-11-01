package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/wager-api/internal/entities"

	"github.com/jackc/pgtype"
)

func Test_Decode(t *testing.T) {
	buf := new(bytes.Buffer)
	wager := &entities.Wager{
		WagerID: pgtype.Int4{Int: 1, Status: pgtype.Present},
		// TotalWagerValue:     20.5,
		// Odds:                40,
		// SellingPercentage:   15.5,
		// CurrentSellingPrice: 16.5,
		// PercentageSold:      23.44,
		// AmountSold:          10,
		// PlacedAt: nil,
	}
	json.NewEncoder(buf).Encode(convertWagerPg2Domain(wager))
	fmt.Println(buf)
}

func Test_RoundNumber(t *testing.T) {
	var number float32 = 12.2345
	newnumber := roundFloat(number)
	fmt.Println(newnumber)

}

func Test_ValidateInputWith2DemimaPoints(t *testing.T) {
	var number float32 = 12.2
	newNumber := number * 100
	// numberInStringFormat := fmt.Sprintf("%f", number)
	// getNumberDecemaPoints := strings.Split(numberInStringFormat, ".")
	// fmt.Println(len(getNumberDecemaPoints[1]))
	// steps := number / step
	if newNumber-float32(int(newNumber)) > 0 {
		fmt.Println("wrong")
	}
}

func Test_Query(t *testing.T) {
	var r http.Request
	r.URL = &url.URL{
		Scheme: "localhost:8080/wagers?page=:1&limit=:12",
	}
	fmt.Println(r.URL.Query().Get("page"))

}
