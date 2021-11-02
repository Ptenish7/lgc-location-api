package api

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/common/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ozonmp/lgc-location-api/internal/repo"

	pb "github.com/ozonmp/lgc-location-api/pkg/lgc-location-api"
)

var (
	_ = promauto.NewCounter(prometheus.CounterOpts{
		Name: "lgc_location_api_location_not_found_total",
		Help: "Total number of locations that were not found",
	})
)

type locationAPI struct {
	pb.UnimplementedLgcLocationApiServiceServer
	repo repo.Repo
}

// NewLocationAPI returns api of lgc-location-api service
func NewLocationAPI(r repo.Repo) pb.LgcLocationApiServiceServer {
	return &locationAPI{repo: r}
}

func (l *locationAPI) CreateLocationV1(
	ctx context.Context,
	req *pb.CreateLocationV1Request,
) (*pb.CreateLocationV1Response, error) {

	log.Debug("LgcLocationApi.CreateLocation: not implemented")

	return nil, status.Error(codes.Internal, "not implemented")
}

func (l *locationAPI) DescribeLocationV1(
	ctx context.Context,
	req *pb.DescribeLocationV1Request,
) (*pb.DescribeLocationV1Response, error) {

	log.Debug("LgcLocationApi.DescribeLocation: not implemented")

	return nil, status.Error(codes.Internal, "not implemented")
}

func (l *locationAPI) ListLocationsV1(
	ctx context.Context,
	req *pb.ListLocationsV1Request,
) (*pb.ListLocationsV1Response, error) {

	log.Debug("LgcLocationApi.ListLocations: not implemented")

	return nil, status.Error(codes.Internal, "not implemented")
}

func (l *locationAPI) RemoveLocationV1(
	ctx context.Context,
	req *pb.RemoveLocationV1Request,
) (*pb.RemoveLocationV1Response, error) {

	log.Debug("LgcLocationApi.RemoveLocation: not implemented")

	return nil, status.Error(codes.Internal, "not implemented")
}
