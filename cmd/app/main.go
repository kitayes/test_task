package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"log/slog"
	"test_task/internal/application"
	delivery "test_task/internal/delivery/http"
	"test_task/internal/repository"
	"test_task/pkg/config"
	"test_task/pkg/logger"
	service "test_task/pkg/services"
)

type Config struct {
	Repo   repository.Config `envPrefix:"REPO_"`
	Logger logger.Config     `envPrefix:"LOGGER_"`
	Http   delivery.Config   `envPrefix:"HTTP_"`
}

// @title           Subscription Service
// @version         1.0
// @description     Service for managing user subscriptions
// @host            localhost:8082
// @BasePath        /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Subscription
func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("failed to load .env file: %s", err.Error())
		return
	}

	cfg := Config{}
	if err := config.ReadEnvConfig(&cfg); err != nil {
		slog.Error("failed to read configs: %s", err.Error())
		return
	}

	log := logger.NewLogger(&cfg.Logger)

	// Репозиторий
	repos := repository.NewRepository(&cfg.Repo, log)
	if err := repos.Init(); err != nil {
		log.Error("failed to init repository: %s", err.Error())
		return
	}

	services := application.NewService(repos)

	handler := delivery.NewHandler(services, &cfg.Http, log)
	if err := handler.Init(); err != nil {
		log.Error("failed to init HTTP handler: %s", err.Error())
		return
	}

	srv := service.NewManager(log)
	srv.AddService(repos, handler)

	ctx := context.Background()
	if err := srv.Run(ctx); err != nil {
		err := errors.Wrap(err, "srv.Run(...) error:")
		log.Error(err.Error())
		return
	}

	log.Info("Subscription service started successfully")
}
