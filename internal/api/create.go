package api

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ozonmp/lgc-location-api/internal/model"
	"github.com/ozonmp/lgc-location-api/internal/pkg/logger"
	pb "github.com/ozonmp/lgc-location-api/pkg/lgc-location-api"
)

func (l *locationAPI) CreateLocationV1(
	ctx context.Context,
	req *pb.CreateLocationV1Request,
) (*pb.CreateLocationV1Response, error) {

	if err := req.Validate(); err != nil {
		logger.ErrorKV(ctx, "CreateLocationV1: invalid argument", "err", err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	locationID, err := l.repo.CreateLocation(ctx, req.Latitude, req.Longitude, req.Title)
	if err != nil {
		logger.ErrorKV(ctx, "CreateLocationV1 failed on repo call", "err", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	event := &model.LocationEvent{
		LocationID: locationID,
		Type:       model.Created,
		Status:     model.Deferred,
		Entity:     &model.Location{Latitude: req.Latitude, Longitude: req.Longitude, Title: req.Title},
	}

	err = l.eventRepo.Add(ctx, event)
	if err != nil {
		logger.ErrorKV(ctx, "failed to add location creation event", "err", err)
	}

	logger.DebugKV(ctx, "CreateLocationV1 succeeded")

	return &pb.CreateLocationV1Response{LocationId: locationID}, nil
}
