package main

import (
	"context"
	"fmt"

	"github/Igo87/crypt/config"
	"github/Igo87/crypt/internal/bot"
	"github/Igo87/crypt/internal/db"
	"github/Igo87/crypt/pkg/api"
	"github/Igo87/crypt/pkg/logger"
	"github/Igo87/crypt/service"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	dbRepo, err := db.NewCurrencyRepository(config.Cfg.GetConnString())
	if err != nil {
		logger.Log.Warn("failed to create db (repo): %v", err)
	}
	defer dbRepo.Close()

	srv := service.NewService(dbRepo)
	if err = srv.Run(ctx, logger.Log); err != nil {
		logger.Log.Warn("failed to create service: %v", err)
	}

	bot, err := bot.Start()
	if err != nil {
		logger.Log.Warn("failed to create bot: %v", err)
	}

	handler := api.NewHandler(srv)
	server := api.New(handler, config.Cfg.GetPort())
	if err = server.Start(); err != nil {
		logger.Log.Warn("failed to start server: %v", err)
	}

	ticker := time.NewTicker(time.Duration(10) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("service stopped by signal")
			return
		case <-ticker.C:
			err := srv.Run(ctx, logger.Log)
			if err != nil {
				logger.Log.Error("failed to run service: %v", err)
			}
			if err = processData(ctx, *srv, *bot); err != nil {
				logger.Log.Error("failed to get data: %v", err)
			}
		}
	}
}

func processData(ctx context.Context, srv service.Service, bot bot.Bot) error {
	dataForToday, err := srv.GetDataByToday(ctx)
	if err != nil {
		return fmt.Errorf("failed to get data by today: %w", err)
	}

	dataForYesterday, err := srv.GetDataByLastDay(ctx)
	if err != nil {
		return fmt.Errorf("failed to get data by last day: %w", err)
	}

	go bot.SendData(ctx, dataForToday, dataForYesterday)
	return nil
}
