package main

import (
	"context"
	"github.com/go-logr/zapr"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jmoiron/sqlx"
	jsoniter "github.com/json-iterator/go"
	"github.com/maksattur/karma8/internal/storage"
	"github.com/pressly/goose"

	_ "github.com/lib/pq"
	"github.com/maksattur/karma8/internal/client"
	"github.com/maksattur/karma8/internal/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	os.Exit(int(run()))
}

func run() statusCode {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP)
	defer cancel()
	zapLog := NewLogger()
	logger := zapr.NewLogger(zapLog)

	cfg, err := NewConfigFromEnv()
	if err != nil {
		logger.Error(err, "failed to load config")
		return StatusCodeFailedLoadConfig
	}

	app := fiber.New(fiber.Config{
		AppName:               "Karma8 1.0.0",
		DisableStartupMessage: true,
		Immutable:             true,
		StrictRouting:         true,
		JSONEncoder:           jsoniter.Marshal,
		JSONDecoder:           jsoniter.Unmarshal,
		ErrorHandler:          server.NewErrorsLoggerMiddleware(fiber.DefaultErrorHandler, logger),
	})

	db, err := sqlx.Connect("postgres", cfg.Postgres.DSN)
	if err != nil {
		logger.Error(err, "failed to connect to PG")
		return StatusCodeFailedConnectPG
	}
	defer db.Close()

	if err = goose.Up(db.DB, "./internal/migrations"); err != nil {
		logger.Error(err, "failed migrations")
		return StatusCodeFailedMigrations
	}

	registerRoutes(
		app.Group("/api/v1"),
		server.NewAPI(storage.NewStorage(db), client.NewClient()),
	)

	app.Use(recover.New())

	go func() {
		<-ctx.Done()
		if err := app.Shutdown(); err != nil {
			logger.Error(err, "shutdown")
			os.Exit(int(StatusCodeFailedStopServer))
		}
	}()

	if err := app.Listen(cfg.Address); err != nil {
		logger.Error(err, "start server")
		return StatusCodeFailedStartServer
	}

	return StatusCodeOK
}
