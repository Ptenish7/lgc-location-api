package log_verbose

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/ozonmp/lgc-location-api/internal/pkg/logger"
)

// UnaryServerInterceptor returns a new unary server interceptor for verbose request and response logging.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, respErr error) {
		resp, respErr = handler(ctx, req)

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return
		}

		verbose := md.Get("Log-Verbose")
		if len(verbose) == 0 || verbose[0] != "true" {
			return
		}

		logger.DebugKV(
			ctx, "got gRPC request",
			"method", info.FullMethod,
			"request", req,
		)

		if respErr == nil {
			logger.DebugKV(
				ctx, "gRPC request succeeded",
				"method", info.FullMethod,
				"response", resp,
			)
		} else {
			logger.DebugKV(
				ctx, "gRPC request failed",
				"method", info.FullMethod,
				"err", respErr,
			)
		}

		return
	}
}
