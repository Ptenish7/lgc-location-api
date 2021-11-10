package api

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ozonmp/lgc-location-api/pkg/lgc-location-api"
)

func (l *locationAPI) ListLocationsV1(
	ctx context.Context,
	req *pb.ListLocationsV1Request,
) (*pb.ListLocationsV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("ListLocationsV1: invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	locations, err := l.repo.ListLocations(ctx, 100, 0)
	if err != nil {
		log.Error().Err(err).Msg("ListLocationsV1 failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	locationsPb := make([]*pb.Location, 0, len(locations))
	for _, location := range locations {
		locationsPb = append(locationsPb, locationToProtobuf(location))
	}

	log.Debug().Msg("ListLocationsV1 succeeded")

	return &pb.ListLocationsV1Response{Locations: locationsPb}, nil
}
