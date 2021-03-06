package api

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ozonmp/lgc-location-api/internal/model"
	"github.com/ozonmp/lgc-location-api/internal/pkg/logger"
	pb "github.com/ozonmp/lgc-location-api/pkg/lgc-location-api"
)

func (l *locationAPI) RemoveLocationV1(
	ctx context.Context,
	req *pb.RemoveLocationV1Request,
) (*pb.RemoveLocationV1Response, error) {

	if err := req.Validate(); err != nil {
		logger.ErrorKV(ctx, "RemoveLocationV1: invalid argument", "err", err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := l.repo.RemoveLocation(ctx, req.LocationId)
	if err != nil {
		logger.ErrorKV(ctx, "RemoveLocationV1 failed on repo call", "err", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	event := &model.LocationEvent{
		LocationID: req.LocationId,
		Type:       model.Removed,
		Status:     model.Deferred,
		Entity:     nil,
	}

	err = l.eventRepo.Add(ctx, event)
	if err != nil {
		logger.ErrorKV(ctx, "failed to add location removal event", "err", err)
	}

	logger.DebugKV(ctx, "RemoveLocationV1 succeeded")

	return &pb.RemoveLocationV1Response{}, nil
}
