package services

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type WagerHandler struct {
	WagerService *WagerService
}

// to handle the request params for handling a paginated request.
func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		limit := r.URL.Query().Get("limit")
		// work both cases: /wagers?page=:page&limit=:limit or /wagers?page=page&limit=limit
		if string(page[0]) == ":" {
			page = page[1:]
		}
		if string(limit[0]) == ":" {
			limit = limit[1:]
		}
		intPage := 0
		// default limit = 10
		intLimit := 10
		var err error
		if page != "" {
			intPage, err = strconv.Atoi(page)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"error": "unable to extract page value",
				})
				return
			}
		}
		if limit != "" {
			intLimit, err = strconv.Atoi(limit)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"error": "unable to extract limit value",
				})
				return
			}
		}
		ctx := context.WithValue(r.Context(), "page", intPage)
		ctx = context.WithValue(ctx, "limit", intLimit)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//	func extractWagerID(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			// page := r.URL.Query().Get("page")
//			// ctx = context.WithValue(ctx, "wager_id", intLimit)
//			// next.ServeHTTP(w, r.WithContext(ctx))
//		})
//	}
func extractWagerID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		var wagerIDInt int
		if wagerID := chi.URLParam(r, "wagerID"); wagerID != "" {
			wagerIDInt, err = strconv.Atoi(wagerID)
			if err != nil {

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"error": "unable to extract limit value",
				})
				return

			}
		}
		ctx := context.WithValue(r.Context(), "wager_id", wagerIDInt)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewWagerHandler(mux *chi.Mux, wagerService *WagerService) {
	handler := &WagerHandler{
		WagerService: wagerService,
	}
	// StripSlashes remove redundant slash in endpoint, example /login/ -> /login
	mux.Use(middleware.StripSlashes)

	mux.Group(func(r chi.Router) {
		// r.Use(handler.authMiddleware)
		r.Post("/wagers", handler.WagerService.PlaceWager)
		r.With(extractWagerID).Post("/buy/{wagerID}", handler.WagerService.BuyWager)
		r.With(paginate).Get("/wagers", handler.WagerService.ListWager)
		// r.Get("/wagers", handler.WagerService.ListWager)

	})
}
