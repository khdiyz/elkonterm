package main

import (
	"context"
	"elkonterm/config"
	"elkonterm/internal/handler"
	"elkonterm/internal/repository"
	"elkonterm/internal/repository/postgres"
	"elkonterm/internal/service"
	"elkonterm/pkg/httpserver"
	"elkonterm/pkg/logger"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

// @title Elkonterm System API
// @version 1.0
// @description API Server for Application
// @host localhost:7070
// @BasePath
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg := config.GetConfig()
	logger := logger.GetLogger()

	db, err := postgres.NewPostgresDB(cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}

	repos := repository.NewRepository(db, logger)
	services := service.NewService(repos, cfg, logger)
	handlers := handler.NewHandler(services, logger)

	srv := new(httpserver.Server)
	go func() {
		if err := srv.Run(cfg.HTTPHost, cfg.HTTPPort, handlers.InitRoutes(cfg)); err != nil {
			logger.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	logger.Info("App started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Warn("App shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logger.Errorf("error occured on db connection close: %s", err.Error())
	}
}
