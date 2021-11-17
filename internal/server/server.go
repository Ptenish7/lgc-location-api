package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"github.com/ozonmp/lgc-location-api/internal/api"
	eventrepo "github.com/ozonmp/lgc-location-api/internal/app/repo"
	"github.com/ozonmp/lgc-location-api/internal/config"
	"github.com/ozonmp/lgc-location-api/internal/pkg/logger"
	"github.com/ozonmp/lgc-location-api/internal/repo"
	"github.com/ozonmp/lgc-location-api/internal/server/middleware/log_level"
	pb "github.com/ozonmp/lgc-location-api/pkg/lgc-location-api"
)

// GrpcServer is gRPC server
type GrpcServer struct {
	db        *sqlx.DB
	batchSize uint
}

// NewGrpcServer returns gRPC server with supporting of batch listing
func NewGrpcServer(db *sqlx.DB, batchSize uint) *GrpcServer {
	return &GrpcServer{
		db:        db,
		batchSize: batchSize,
	}
}

// Start method runs server
func (s *GrpcServer) Start(ctx context.Context, cfg *config.Config) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gatewayAddr := fmt.Sprintf("%s:%v", cfg.Rest.Host, cfg.Rest.Port)
	grpcAddr := fmt.Sprintf("%s:%v", cfg.Grpc.Host, cfg.Grpc.Port)
	metricsAddr := fmt.Sprintf("%s:%v", cfg.Metrics.Host, cfg.Metrics.Port)

	gatewayServer := createGatewayServer(ctx, grpcAddr, gatewayAddr)

	go func() {
		logger.InfoKV(ctx, fmt.Sprintf("gateway server is running on %s", gatewayAddr))
		if err := gatewayServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorKV(ctx, "failed to run gateway server", "err", err)
			cancel()
		}
	}()

	metricsServer := createMetricsServer(cfg)

	go func() {
		logger.InfoKV(ctx, fmt.Sprintf("metrics server is running on %s", metricsAddr))
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorKV(ctx, "failed to run metrics server", "err", err)
			cancel()
		}
	}()

	isReady := &atomic.Value{}
	isReady.Store(false)

	statusServer := createStatusServer(ctx, cfg, isReady)

	go func() {
		statusAddr := fmt.Sprintf("%s:%v", cfg.Status.Host, cfg.Status.Port)
		logger.InfoKV(ctx, fmt.Sprintf("status server is running on %s", statusAddr))
		if err := statusServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorKV(ctx, "failed to run status server", "err", err)
		}
	}()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer func() {
		if err := l.Close(); err != nil {
			logger.ErrorKV(ctx, "failed to close listener")
		}
	}()

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: time.Duration(cfg.Grpc.MaxConnectionIdle) * time.Minute,
			Timeout:           time.Duration(cfg.Grpc.Timeout) * time.Second,
			MaxConnectionAge:  time.Duration(cfg.Grpc.MaxConnectionAge) * time.Minute,
			Time:              time.Duration(cfg.Grpc.Timeout) * time.Minute,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_opentracing.UnaryServerInterceptor(),
			grpcrecovery.UnaryServerInterceptor(),
			log_level.UnaryServerInterceptor(),
		)),
	)

	r := repo.NewRepo(s.db, s.batchSize)
	er := eventrepo.NewEventRepo(s.db)

	pb.RegisterLgcLocationApiServiceServer(grpcServer, api.NewLocationAPI(r, er))
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(grpcServer)

	go func() {
		logger.InfoKV(ctx, fmt.Sprintf("gRPC server is listening on %s", grpcAddr))
		if err := grpcServer.Serve(l); err != nil {
			logger.FatalKV(ctx, "failed to run gRPC server", "err", err)
		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		isReady.Store(true)
		logger.InfoKV(ctx, "the service is ready to accept requests")
	}()

	if cfg.Project.Debug {
		reflection.Register(grpcServer)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		logger.InfoKV(ctx, fmt.Sprintf("signal.Notify: %v", v))
	case done := <-ctx.Done():
		logger.InfoKV(ctx, fmt.Sprintf("ctx.Done: %v", done))
	}

	isReady.Store(false)

	if err := gatewayServer.Shutdown(ctx); err != nil {
		logger.ErrorKV(ctx, "gatewayServer.Shutdown failed", "err", err)
	} else {
		logger.InfoKV(ctx, "gatewayServer shut down correctly")
	}

	if err := statusServer.Shutdown(ctx); err != nil {
		logger.ErrorKV(ctx, "statusServer.Shutdown failed", "err", err)
	} else {
		logger.InfoKV(ctx, "statusServer shut down correctly")
	}

	if err := metricsServer.Shutdown(ctx); err != nil {
		logger.ErrorKV(ctx, "metricsServer.Shutdown failed", "err", err)
	} else {
		logger.InfoKV(ctx, "metricsServer shut down correctly")
	}

	grpcServer.GracefulStop()
	logger.InfoKV(ctx, "grpcServer shut down correctly")

	return nil
}
