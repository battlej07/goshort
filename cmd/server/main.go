package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/battlej07/goshort/internal/config"
	"github.com/battlej07/goshort/internal/store"
)

type application struct {
	logger   *slog.Logger
	urlStore *store.ShortURLStore
	config   *config.Config
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Error loading config", "error", err)
		return
	}

	dbCtx, cancelDB := context.WithTimeout(context.Background(), 5*time.Second)
	urlStore, err := store.NewShortURLStore(dbCtx, cfg.DatabaseURL)
	cancelDB()
	if err != nil {
		logger.Error("error connecting to database", "error", err)
		return
	}
	defer urlStore.Close()

	schemaCtx, cancelSchema := context.WithTimeout(context.Background(), 5*time.Second)
	err = urlStore.Init(schemaCtx)
	cancelSchema()
	if err != nil {
		logger.Error("error initializing database schema", "error", err)
		return
	}

	app := &application{
		logger:   logger,
		urlStore: urlStore,
		config:   cfg,
	}

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		logger.Info("starting server", "port", cfg.Port)
		if serveErr := srv.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			logger.Error("server failed", "error", serveErr)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("error shutting down server", "error", err)
	}
}
