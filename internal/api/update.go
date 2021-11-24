package api

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ozonmp/lgc-location-api/internal/model"
	"github.com/ozonmp/lgc-location-api/internal/pkg/logger"
	pb "github.com/ozonmp/lgc-location-api/pkg/lgc-location-api"
)

func (l *locationAPI) UpdateLocationV1(
	ctx context.Context,
	req *pb.UpdateLocationV1Request,
) (*pb.UpdateLocationV1Response, error) {

	if err := req.Validate(); err != nil {
		logger.ErrorKV(ctx, "UpdateLocationV1: invalid argument", "err", err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	location := &model.Location{
		ID:        req.GetLocationId(),
		Latitude:  req.GetLatitude(),
		Longitude: req.GetLongitude(),
		Title:     req.GetTitle(),
	}

	err := l.repo.UpdateLocation(ctx, location)
	if err != nil {
		logger.ErrorKV(ctx, "UpdateLocationV1 failed on repo call", "err", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	event := &model.LocationEvent{
		LocationID: req.LocationId,
		Type:       model.Updated,
		TypeExtra:  model.FullEventTypeExtra(),
		Status:     model.Deferred,
		Entity:     location,
	}

	err = l.eventRepo.Add(ctx, event)
	if err != nil {
		logger.ErrorKV(ctx, "failed to add location update event", "err", err)
	}

	logger.DebugKV(ctx, "UpdateLocationV1 succeeded")

	return &pb.UpdateLocationV1Response{}, nil
}
