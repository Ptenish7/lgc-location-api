package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	gelf "github.com/snovichkov/zap-gelf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"

	"github.com/ozonmp/lgc-location-api/internal/config"
	"github.com/ozonmp/lgc-location-api/internal/database"
	"github.com/ozonmp/lgc-location-api/internal/pkg/logger"
	eventrepo "github.com/ozonmp/lgc-location-api/internal/retranslator/repo"
	"github.com/ozonmp/lgc-location-api/internal/retranslator/retranslator"
	"github.com/ozonmp/lgc-location-api/internal/retranslator/sender"
	"github.com/ozonmp/lgc-location-api/internal/server"
	"github.com/ozonmp/lgc-location-api/internal/tracer"
)

var (
	batchSize uint = 2
)

func main() {
	ctx := context.Background()

	if err := config.ReadConfigYML("config.yml"); err != nil {
		logger.FatalKV(ctx, "failed to init configuration", "err", err)
	}
	cfg := config.GetConfigInstance()

	migration := flag.Bool("migration", true, "Defines the migration start option")
	flag.Parse()

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

	db, err := database.NewPostgres(ctx, dsn, cfg.Database.Driver)
	if err != nil {
		logger.FatalKV(ctx, "failed to init postgres", "err", err)
	}
	defer func() {
		if clErr := db.Close(); clErr != nil {
			logger.ErrorKV(ctx, "failed to close database connection", "err", clErr)
		}
	}()

	if *migration {
		if err = goose.Up(db.DB, cfg.Database.Migrations); err != nil {
			logger.ErrorKV(ctx, "migration failed", "err", err)

			return
		}
	}

	tracing, err := tracer.NewTracer(ctx, &cfg)
	if err != nil {
		logger.ErrorKV(ctx, "failed to init tracing", "err", err)

		return
	}
	defer func() {
		if clErr := tracing.Close(); clErr != nil {
			logger.ErrorKV(ctx, "failed to close tracer", "err", clErr)
		}
	}()

	rtConfig, err := configRetranslator(db, &cfg.Kafka)
	if err != nil {
		logger.ErrorKV(ctx, "failed to config retranslator", "err", err)

		return
	}
	rt := retranslator.NewRetranslator(*rtConfig)
	go rt.Start()
	defer rt.Close()

	if err := server.NewGrpcServer(db, batchSize).Start(ctx, &cfg); err != nil {
		logger.ErrorKV(ctx, "failed to create gRPC server", "err", err)

		return
	}
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

func configRetranslator(db *sqlx.DB, kafkaConfig *config.Kafka) (*retranslator.Config, error) {
	eventSender, err := sender.NewEventSender(kafkaConfig.Brokers)
	if err != nil {
		return nil, err
	}

	return &retranslator.Config{
		ChannelSize:     16,
		ConsumerCount:   2,
		ConsumerSize:    4,
		ConsumerTimeout: 10 * time.Second,
		ProducerCount:   2,
		WorkerCount:     2,
		BatchSize:       4,
		Repo:            eventrepo.NewEventRepo(db),
		Sender:          eventSender,
	}, nil
}
