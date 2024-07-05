package main

import (
	"log/slog"
	"os"

	_ "github.com/njslxve/time-tracker-service/docs"
	"github.com/njslxve/time-tracker-service/internal/config"
	"github.com/njslxve/time-tracker-service/internal/server"
	"github.com/njslxve/time-tracker-service/internal/service"
	"github.com/njslxve/time-tracker-service/internal/transport/api"
	"github.com/njslxve/time-tracker-service/internal/transport/storage"
	"github.com/njslxve/time-tracker-service/pkg/client/postgres"
	"github.com/njslxve/time-tracker-service/pkg/logger"
)

// @title Time Tracker API
// @version 1.0

// @host localhost:8080
// @BasePath /
func main() {
	logger := logger.New()
	slog.SetDefault(logger)

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Debug("config error: ",
			slog.String("error", err.Error()))
		slog.Error("failed to load config",
			slog.String("error", err.Error()))
		os.Exit(1)
	}

	client, err := postgres.NewClient(cfg)
	if err != nil {
		slog.Debug("db error: ",
			slog.String("error", err.Error()))
		slog.Error("failed to connect to db",
			slog.String("error", err.Error()))
		os.Exit(1)
	}

	storage := storage.New(logger, client)
	api := api.New(logger, cfg)

	service := service.New(cfg, logger, storage, api)

	server := server.New(cfg, logger, service)

	server.Start()
}
