package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/aybancaughtin4k/gymops/backend/internal/data"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

type config struct {
	port         int
	env          string
	dsn          string
	authTokenKey string
}

type application struct {
	logger *slog.Logger
	config config
	dbpool *pgxpool.Pool
	models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "The port of server that will use")
	flag.StringVar(&cfg.env, "environment", "development", "Server environment (production|development)")
	flag.StringVar(&cfg.dsn, "dsn", os.Getenv("GOOSE_DBSTRING"), "Database dsn")
	flag.StringVar(&cfg.authTokenKey, "token-key", os.Getenv("AUTH_TOKEN_KEY"), "Authentication token key secret")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	dbpool, err := initDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("initializing database connection")

	defer dbpool.Close()

	app := &application{
		logger: logger,
		config: cfg,
		dbpool: dbpool,
		models: data.NewModels(dbpool),
	}

	err = app.serve()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}

func initDB(cfg config) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), cfg.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = dbpool.Ping(ctx)
	if err != nil {
		dbpool.Close()
		return nil, err
	}

	return dbpool, nil
}
