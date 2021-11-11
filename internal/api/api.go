package api

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/ozonmp/lgc-location-api/internal/app/repo"
	"github.com/ozonmp/lgc-location-api/internal/model"
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
	repo      repo.Repo
	eventRepo eventrepo.EventRepo
}

// NewLocationAPI returns api of lgc-location-api service
func NewLocationAPI(r repo.Repo, er eventrepo.EventRepo) pb.LgcLocationApiServiceServer {
	return &locationAPI{repo: r, eventRepo: er}
}

func locationToProtobuf(l *model.Location) *pb.Location {
	return &pb.Location{
		Id:        l.ID,
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
		Title:     l.Title,
	}
}
