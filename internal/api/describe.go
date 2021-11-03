package api

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ozonmp/lgc-location-api/pkg/lgc-location-api"
)

func (l *locationAPI) DescribeLocationV1(
	ctx context.Context,
	req *pb.DescribeLocationV1Request,
) (*pb.DescribeLocationV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("DescribeLocationV1: invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	location, err := l.repo.DescribeLocation(ctx, req.LocationId)
	if err != nil {
		log.Error().Err(err).Msg("DescribeLocationV1 failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Msg("DescribeLocationV1 succeeded")

	return &pb.DescribeLocationV1Response{Location: locationToProtobuf(location)}, nil
}
