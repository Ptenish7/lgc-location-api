package api

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ozonmp/lgc-location-api/pkg/lgc-location-api"
)

func (l *locationAPI) CreateLocationV1(
	ctx context.Context,
	req *pb.CreateLocationV1Request,
) (*pb.CreateLocationV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("CreateLocationV1: invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	locationID, err := l.repo.CreateLocation(ctx, req.Latitude, req.Longitude, req.Title)
	if err != nil {
		log.Error().Err(err).Msg("CreateLocationV1 failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Msg("CreateLocationV1 succeeded")

	return &pb.CreateLocationV1Response{LocationId: locationID}, nil
}