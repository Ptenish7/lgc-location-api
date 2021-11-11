package eventrepo

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"

	"github.com/ozonmp/lgc-location-api/internal/model"
)

type EventRepoTestSuite struct {
	suite.Suite
	r    EventRepo
	mock sqlmock.Sqlmock
}

func (s *EventRepoTestSuite) SetupSuite() {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		s.T().Fatal(err)
	}

	s.mock = mock
	s.r = NewEventRepo(sqlx.NewDb(mockDB, "sqlmock"))
}

func (s *EventRepoTestSuite) TestLock() {
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	s.mock.
		ExpectQuery(regexp.QuoteMeta(`UPDATE locations_events SET status = $1, updated_at = $2 WHERE id IN ( SELECT id FROM locations_events WHERE status = $3 LIMIT 1 ) RETURNING *`)).
		WithArgs("Processed", "now()", "Deferred").
		WillReturnRows(rows)

	_, err := s.r.Lock(context.Background(), 1)
	s.Nil(err)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *EventRepoTestSuite) TestUnlock() {
	s.mock.
		ExpectExec(regexp.QuoteMeta(`UPDATE locations_events SET status = $1, updated_at = $2 WHERE id IN ($3)`)).
		WithArgs("Deferred", "now()", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := s.r.Unlock(context.Background(), []uint64{1})
	s.Nil(err)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *EventRepoTestSuite) TestAdd() {
	s.mock.
		ExpectExec(regexp.QuoteMeta(`INSERT INTO locations_events (location_id,type,type_extra,status,payload) VALUES ($1,$2,$3,$4,$5)`)).
		WithArgs(1, "Removed", 0, "Deferred", []byte("{}")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.r.Add(context.Background(), &model.LocationEvent{LocationID: 1, Type: model.Removed, Status: model.Deferred, Entity: nil})
	s.Nil(err)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *EventRepoTestSuite) TestRemove() {
	s.mock.
		ExpectExec(regexp.QuoteMeta(`DELETE FROM locations_events WHERE id IN ($1)`)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := s.r.Remove(context.Background(), []uint64{1})
	s.Nil(err)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEventRepo(t *testing.T) {
	suite.Run(t, new(EventRepoTestSuite))
}
