package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	gelf "github.com/snovichkov/zap-gelf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ozonmp/lgc-location-api/internal/database"
	"github.com/ozonmp/lgc-location-api/internal/pkg/logger"
	"github.com/ozonmp/lgc-location-api/internal/retranslator/config"
	eventrepo "github.com/ozonmp/lgc-location-api/internal/retranslator/repo"
	"github.com/ozonmp/lgc-location-api/internal/retranslator/retranslator"
	"github.com/ozonmp/lgc-location-api/internal/retranslator/sender"
)

func main() {
	ctx := context.Background()

	if err := config.ReadConfigYML("config-rt.yaml"); err != nil {
		logger.FatalKV(ctx, "failed to init retranslator configuration", "err", err)
	}
	cfg := config.GetConfigInstance()

	syncLogger := initLogger(ctx, cfg)
	defer syncLogger()

	logger.InfoKV(
		ctx, "starting service",
		"name", cfg.Project.Name,
		"version", cfg.Project.Version,
		"commitHash", cfg.Project.CommitHash,
		"debug", cfg.Project.Debug,
		"environment", cfg.Project.Environment,
	)

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)

	db, err := database.NewPostgres(ctx, dsn, cfg.Database.Driver, cfg.Database.MaxRetry)
	if err != nil {
		logger.FatalKV(ctx, "failed to init postgres", "err", err)
	}
	defer func() {
		if clErr := db.Close(); clErr != nil {
			logger.ErrorKV(ctx, "failed to close database connection", "err", clErr)
		}
	}()

	eventSender, err := sender.NewEventSender(cfg.Kafka.Brokers, cfg.Kafka.MaxRetry)
	if err != nil {
		logger.FatalKV(ctx, "failed to init event sender", "err", err)
	}

	eventRepo := eventrepo.NewEventRepo(db)

	rt := retranslator.NewRetranslator(&cfg.Retranslator, eventRepo, eventSender)
	go rt.Start()
	defer rt.Close()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

func initLogger(ctx context.Context, cfg config.Config) (syncFn func()) {
	loggingLevel := zap.InfoLevel
	if cfg.Project.Debug {
		loggingLevel = zap.DebugLevel
	}

	consoleCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		os.Stderr,
		zap.NewAtomicLevelAt(loggingLevel),
	)

	gelfCore, err := gelf.NewCore(
		gelf.Addr(cfg.Telemetry.GraylogPath),
		gelf.Level(loggingLevel),
	)
	if err != nil {
		logger.FatalKV(ctx, "gelf.NewCore() error", "err", err)
	}

	notSugaredLogger := zap.New(zapcore.NewTee(consoleCore, gelfCore))

	sugaredLogger := notSugaredLogger.Sugar()
	logger.SetLogger(sugaredLogger.With("service", cfg.Project.Name))

	return func() {
		if syncErr := notSugaredLogger.Sync(); syncErr != nil {
			logger.ErrorKV(ctx, "failed to sync logger", "err", syncErr)
		}
	}
}
