package api

import (
	"github.com/ozonmp/lgc-location-api/internal/model"
	"github.com/ozonmp/lgc-location-api/internal/repo"
	"github.com/ozonmp/lgc-location-api/internal/retranslator/repo"

	pb "github.com/ozonmp/lgc-location-api/pkg/lgc-location-api"
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
