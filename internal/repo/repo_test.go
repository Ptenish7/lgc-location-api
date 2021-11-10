package repo

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type RepoTestSuite struct {
	suite.Suite
	r    Repo
	mock sqlmock.Sqlmock
}

func (s *RepoTestSuite) SetupSuite() {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		s.T().Fatal(err)
	}

	s.mock = mock
	s.r = NewRepo(sqlx.NewDb(mockDB, "sqlmock"), 1)
}

func (s *RepoTestSuite) TestCreateLocation() {
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	s.mock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO locations (latitude,longitude,title) VALUES ($1,$2,$3)`)).
		WithArgs(10.0, 20.0, "L1").
		WillReturnRows(rows)

	_, err := s.r.CreateLocation(context.Background(), 10, 20, "L1")
	s.Nil(err)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *RepoTestSuite) TestDescribeLocation() {
	rows := sqlmock.
		NewRows([]string{"id", "latitude", "longitude", "title", "removed", "created_at", "updated_at"}).
		AddRow(1, 10.0, 20.0, "L1", false, time.Now(), time.Now())

	s.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM locations WHERE id = $1 AND removed = $2`)).
		WithArgs(1, false).
		WillReturnRows(rows)

	_, err := s.r.DescribeLocation(context.Background(), 1)
	s.Nil(err)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *RepoTestSuite) TestListLocation() {
	rows := sqlmock.
		NewRows([]string{"id", "latitude", "longitude", "title", "removed", "created_at", "updated_at"}).
		AddRow(1, 10.0, 20.0, "L1", false, time.Now(), time.Now())

	s.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM locations WHERE removed = $1 ORDER BY id LIMIT 1 OFFSET 0`)).
		WithArgs(false).
		WillReturnRows(rows)

	_, err := s.r.ListLocations(context.Background(), 1, 0)
	s.Nil(err)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *RepoTestSuite) TestRemoveLocation() {
	s.mock.
		ExpectExec(regexp.QuoteMeta(`UPDATE locations SET removed = $1, updated_at = $2 WHERE id = $3 AND removed = $4`)).
		WithArgs(true, "now()", 1, false).
		WillReturnResult(sqlmock.NewResult(0, 1))

	_, err := s.r.RemoveLocation(context.Background(), 1)
	s.Nil(err)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRepo(t *testing.T) {
	suite.Run(t, new(RepoTestSuite))
}
