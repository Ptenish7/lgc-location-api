package api

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ozonmp/lgc-location-api/pkg/lgc-location-api"
)

func (l *locationAPI) RemoveLocationV1(
	ctx context.Context,
	req *pb.RemoveLocationV1Request,
) (*pb.RemoveLocationV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("RemoveLocationV1: invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	found, err := l.repo.RemoveLocation(ctx, req.LocationId)
	if err != nil {
		log.Error().Err(err).Msg("RemoveLocationV1 failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Msg("RemoveLocationV1 succeeded")

	return &pb.RemoveLocationV1Response{Found: found}, nil
}
