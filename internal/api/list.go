package api

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ozonmp/lgc-location-api/internal/pkg/logger"
	pb "github.com/ozonmp/lgc-location-api/pkg/lgc-location-api"
)

func (l *locationAPI) ListLocationsV1(
	ctx context.Context,
	req *pb.ListLocationsV1Request,
) (*pb.ListLocationsV1Response, error) {

	if err := req.Validate(); err != nil {
		logger.ErrorKV(ctx, "ListLocationsV1: invalid argument", "err", err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	locations, err := l.repo.ListLocations(ctx, 100, 0)
	if err != nil {
		logger.ErrorKV(ctx, "ListLocationsV1 failed on repo call")

		return nil, status.Error(codes.Internal, err.Error())
	}

	locationsPb := make([]*pb.Location, 0, len(locations))
	for _, location := range locations {
		locationsPb = append(locationsPb, locationToProtobuf(location))
	}

	logger.DebugKV(ctx, "ListLocationsV1 succeeded")

	return &pb.ListLocationsV1Response{Locations: locationsPb}, nil
}
