package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/battlej07/goshort/internal/config"
)

type application struct {
	logger *slog.Logger
	db     map[string]string
	config *config.Config
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db := map[string]string{}

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Error loading config", "error", err)
		return
	}

	app := &application{
		logger: logger,
		db:     db,
		config: cfg,
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
		err = srv.ListenAndServe()
		logger.Error(err.Error())
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
}
