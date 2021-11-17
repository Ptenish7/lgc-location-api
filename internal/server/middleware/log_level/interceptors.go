package log_level

import (
	"context"
	"strings"

	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/ozonmp/lgc-location-api/internal/pkg/logger"
)

// UnaryServerInterceptor returns a new unary server interceptor for redefining log level.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			levels := md.Get("Log-Level")
			if len(levels) > 0 {
				logger.InfoKV(ctx, "got log level", "levels", levels)
				if parsedLevel, ok := parseLevel(levels[0]); ok {
					newLogger := logger.CloneWithLevel(ctx, parsedLevel)
					ctx = logger.AttachLogger(ctx, newLogger)
				}
			}
		}
		return handler(ctx, req)
	}
}

func parseLevel(str string) (zapcore.Level, bool) {
	switch strings.ToLower(str) {
	case "debug":
		return zapcore.DebugLevel, true
	case "info":
		return zapcore.InfoLevel, true
	case "warn":
		return zapcore.WarnLevel, true
	case "error":
		return zapcore.ErrorLevel, true
	default:
		return zapcore.DebugLevel, false
	}
}
