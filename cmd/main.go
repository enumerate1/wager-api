package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/wager-api/internal/repositories"
	"github.com/wager-api/internal/services"
	"github.com/wager-api/libs/configs"
	"github.com/wager-api/libs/database"
	"github.com/wager-api/libs/logs"
	"github.com/wager-api/libs/mux"

	"go.uber.org/zap"
)

func main() {
	var err error
	ctx := context.Background()
	// ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	// defer stop()
	// load the config, load the file default, and load the environment variable
	cfg := configs.LoadConfigFile("./configs/config.yaml")
	// configs.LoadConfigEnv(&cfg)

	logs.Logger, err = logs.InitWithOption(cfg.LogLevel, cfg.Service)
	if err != nil {
		log.Println("can't setup zap log", err)
	}
	zap.ReplaceGlobals(logs.Logger.Desugar())
	defer logs.Logger.Sync()

	// connect DB

	pool := database.NewConnectionPool(ctx, logs.Logger.Desugar(), cfg.Postgres)

	// wagerService := services
	wagerService := &services.WagerService{
		DB:           pool,
		WagerRepo:    &repositories.WagerRepo{},
		PurchaseRepo: &repositories.PurchaseRepo{},
	}

	mux := mux.InitWithLogger(logs.Logger.Desugar())
	services.NewWagerHandler(mux, wagerService)
	// logging.Logger.Infof("Listening at %s", cfg.Address)
	err = http.ListenAndServe(cfg.Address, mux)
	if err != nil {
		fmt.Println("====err", err)
	}
}
