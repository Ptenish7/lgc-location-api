package api

import (
	"context"
	"database/sql"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ozonmp/lgc-location-api/internal/metrics"
	"github.com/ozonmp/lgc-location-api/internal/pkg/logger"
	pb "github.com/ozonmp/lgc-location-api/pkg/lgc-location-api"
)

func (l *locationAPI) DescribeLocationV1(
	ctx context.Context,
	req *pb.DescribeLocationV1Request,
) (*pb.DescribeLocationV1Response, error) {

	if err := req.Validate(); err != nil {
		logger.ErrorKV(ctx, "DescribeLocationV1: invalid argument", "err", err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	location, err := l.repo.DescribeLocation(ctx, req.LocationId)
	if err != nil {
		logger.ErrorKV(ctx, "DescribeLocationV1 failed on repo call", "err", err)

		if errors.Is(err, sql.ErrNoRows) {
			metrics.IncLocationNotFoundCounter()
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.DebugKV(ctx, "DescribeLocationV1 succeeded")

	return &pb.DescribeLocationV1Response{Location: locationToProtobuf(location)}, nil
}
