package api

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"github.com/ozonmp/lgc-location-api/internal/mocks"
	pb "github.com/ozonmp/lgc-location-api/pkg/lgc-location-api"
)

const bufSize = 1024 * 1024

type LocationAPITestSuite struct {
	suite.Suite
	listener *bufconn.Listener
	server   *grpc.Server
	conn     *grpc.ClientConn
	client   pb.LgcLocationApiServiceClient
}

func (s *LocationAPITestSuite) bufDialer(context.Context, string) (net.Conn, error) {
	return s.listener.Dial()
}

func (s *LocationAPITestSuite) SetupSuite() {
	s.listener = bufconn.Listen(bufSize)
	s.server = grpc.NewServer()

	ctrl := gomock.NewController(s.T())
	repo := mocks.NewMockRepo(ctrl)

	pb.RegisterLgcLocationApiServiceServer(s.server, NewLocationAPI(repo))
	go func() {
		if err := s.server.Serve(s.listener); err != nil {
			log.Fatalf("s.server exited with error: %v", err)
		}
	}()

	ctx := context.Background()
	var err error
	s.conn, err = grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(s.bufDialer), grpc.WithInsecure())
	if err != nil {
		s.T().Fatalf("failed to dial bufnet: %v", err)
	}

	s.client = pb.NewLgcLocationApiServiceClient(s.conn)
}

func (s *LocationAPITestSuite) TearDownSuite() {
	err := s.conn.Close()
	if err != nil {
		log.Panic(err)
	}
	s.server.Stop()
}

func (s *LocationAPITestSuite) TestCreateLocation_LatitudeValidation() {
	req := &pb.CreateLocationV1Request{
		Latitude:  -120,
		Longitude: 10,
		Title:     "L1",
	}

	resp, err := s.client.CreateLocationV1(context.Background(), req)
	s.Nil(resp)
	s.NotNil(err)

	st, _ := status.FromError(err)
	s.Equal(codes.InvalidArgument, st.Code())
	s.Equal("invalid CreateLocationV1Request.Latitude: value must be inside range [-90, 90]", st.Message())
}

func (s *LocationAPITestSuite) TestCreateLocation_LongitudeValidation() {
	req := &pb.CreateLocationV1Request{
		Latitude:  -10,
		Longitude: 220,
		Title:     "L1",
	}

	resp, err := s.client.CreateLocationV1(context.Background(), req)
	s.Nil(resp)
	s.NotNil(err)

	st, _ := status.FromError(err)
	s.Equal(codes.InvalidArgument, st.Code())
	s.Equal("invalid CreateLocationV1Request.Longitude: value must be inside range [-180, 180]", st.Message())
}

func (s *LocationAPITestSuite) TestCreateLocation_TitleValidation() {
	req := &pb.CreateLocationV1Request{
		Latitude:  -10,
		Longitude: 10,
		Title:     "",
	}

	resp, err := s.client.CreateLocationV1(context.Background(), req)
	s.Nil(resp)
	s.NotNil(err)

	st, _ := status.FromError(err)
	s.Equal(codes.InvalidArgument, st.Code())
	s.Equal("invalid CreateLocationV1Request.Title: value length must be at least 1 runes", st.Message())
}

func (s *LocationAPITestSuite) TestDescribeLocation_LocationIdValidation() {
	req := &pb.DescribeLocationV1Request{LocationId: 0}

	resp, err := s.client.DescribeLocationV1(context.Background(), req)
	s.Nil(resp)
	s.NotNil(err)

	st, _ := status.FromError(err)
	s.Equal(codes.InvalidArgument, st.Code())
	s.Equal("invalid DescribeLocationV1Request.LocationId: value must be greater than 0", st.Message())
}

func (s *LocationAPITestSuite) TestRemoveLocation_LocationIdValidation() {
	req := &pb.RemoveLocationV1Request{LocationId: 0}

	resp, err := s.client.RemoveLocationV1(context.Background(), req)
	s.Nil(resp)
	s.NotNil(err)

	st, _ := status.FromError(err)
	s.Equal(codes.InvalidArgument, st.Code())
	s.Equal("invalid RemoveLocationV1Request.LocationId: value must be greater than 0", st.Message())
}

func TestLocationAPI(t *testing.T) {
	suite.Run(t, new(LocationAPITestSuite))
}
