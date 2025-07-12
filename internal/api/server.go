package api

import (
	"fmt"
	"log/slog"

	"github.com/ffss92/webhookd/internal/config"
	"github.com/ffss92/webhookd/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServerConfig struct {
	Config  *config.Config
	DevMode bool
	Logger  *slog.Logger
	Pool    *pgxpool.Pool
}

type Server struct {
	devMode bool
	cfg     *config.Config
	logger  *slog.Logger
	pool    *pgxpool.Pool
	store   *database.Store
}

func NewServer(scfg ServerConfig) (*Server, error) {
	if scfg.Config == nil {
		return nil, fmt.Errorf("missing config in server config")
	}
	if scfg.Pool == nil {
		return nil, fmt.Errorf("missing db pool in server config")
	}

	return &Server{
		devMode: scfg.DevMode,
		logger:  scfg.Logger,
		cfg:     scfg.Config,
		pool:    scfg.Pool,
		store:   database.New(scfg.Pool),
	}, nil
}
