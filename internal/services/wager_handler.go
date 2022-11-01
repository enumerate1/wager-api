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
func paginateMiddleware(next http.Handler) http.Handler {
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

func extractWagerIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		var wagerIDInt int
		if wagerID := chi.URLParam(r, "wagerID"); wagerID != "" {
			wagerIDInt, err = strconv.Atoi(wagerID)
			if err != nil {
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
func setContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func NewWagerHandler(mux *chi.Mux, wagerService *WagerService) {
	handler := &WagerHandler{
		WagerService: wagerService,
	}
	// StripSlashes remove redundant slash in endpoint, example /login/ -> /login
	mux.Use(middleware.StripSlashes)
	mux.Use(setContentTypeMiddleware)

	mux.Group(func(r chi.Router) {
		r.Post("/wagers", handler.WagerService.PlaceWager)
		r.With(extractWagerIDMiddleware).Post("/buy/{wagerID}", handler.WagerService.BuyWager)
		r.With(paginateMiddleware).Get("/wagers", handler.WagerService.ListWager)
	})
}
