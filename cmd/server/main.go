package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/ffss92/webhookd/internal/api"
	"github.com/ffss92/webhookd/internal/config"
	"github.com/ffss92/webhookd/internal/postgres"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

var (
	devMode bool
)

func run() error {
	flag.BoolVar(&devMode, "dev", false, "Sets the application in dev mode")
	flag.Parse()

	ctx := context.Background()

	cfg, err := config.NewFromEnv()
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	pool, err := postgres.New(ctx, cfg.DBConn())
	if err != nil {
		return err
	}
	defer pool.Close()

	apisrv, err := api.NewServer(api.ServerConfig{
		Config:  cfg,
		DevMode: devMode,
		Pool:    pool,
		Logger:  logger,
	})
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:     cfg.Addr(),
		Handler:  apisrv.Routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", slog.String("addr", srv.Addr), slog.Bool("dev", devMode))
	return srv.ListenAndServe()
}
