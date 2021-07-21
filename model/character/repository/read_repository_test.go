package repository_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	redismock "github.com/elliotchance/redismock/v8"
	redis "github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/hezbymuhammad/golang-marvel-demo/domain"
	"github.com/hezbymuhammad/golang-marvel-demo/model/character/repository"
)

var record domain.Character

type CharacterReadRepositoryTestSuite struct {
	suite.Suite
	mock *redismock.ClientMock
	repo domain.CharacterReadRepository
}

func TestCharacterReadRepository(t *testing.T) {
	suite.Run(t, new(CharacterReadRepositoryTestSuite))
}

func (s *CharacterReadRepositoryTestSuite) SetupTest() {
	mr, err := miniredis.Run()
	if err != nil {
		s.T().Fatalf("Error: '%s'", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	s.mock = redismock.NewNiceMock(client)
	s.repo = repository.NewCharacterReadRepository(s.mock)
}

func (s *CharacterReadRepositoryTestSuite) TestNilFetch() {
	s.mock.On("Exists", mock.Anything, mock.Anything).Return(redis.NewIntResult(1, nil))
	s.mock.On("Get", mock.Anything, "marvel-characters-page-1").Return(redis.NewStringResult("", nil))

	_, err := s.repo.Fetch(context.Background(), 1)
	s.Assert().Equal(err, domain.ErrNotFound)
}

func (s *CharacterReadRepositoryTestSuite) TestFailedJSONFetch() {
	s.mock.On("Exists", mock.Anything, mock.Anything).Return(redis.NewIntResult(1, nil))
	s.mock.On("Get", mock.Anything, "marvel-characters-page-2").Return(redis.NewStringResult("val", nil))

	_, err := s.repo.Fetch(context.Background(), 2)
	s.Assert().Equal(err, domain.ErrInternalServerError)
}

func (s *CharacterReadRepositoryTestSuite) TestFailedEmptyKeyFetch() {
	s.mock.On("Exists", mock.Anything, []string{"marvel-characters-page-2"}).Return(redis.NewIntResult(0, nil))

	_, err := s.repo.Fetch(context.Background(), 2)
	s.Assert().Equal(err, domain.ErrCacheKeyEmpty)
}

func (s *CharacterReadRepositoryTestSuite) TestFailedErrorEmptyKeyFetch() {
	s.mock.On("Exists", mock.Anything, []string{"marvel-characters-page-2"}).Return(redis.NewIntResult(0, errors.New("fail")))

	_, err := s.repo.Fetch(context.Background(), 2)
	s.Assert().Equal(err, domain.ErrInternalServerError)
}

func (s *CharacterReadRepositoryTestSuite) TestSuccessPageNilFetch() {
	IDs := []int{1, 2, 3}
	json_data, err := json.Marshal(IDs)
	if err != nil {
		s.T().Fatalf("Error: '%s'", err)
	}

	s.mock.On("Get", mock.Anything, "marvel-characters-page-1").Return(redis.NewStringResult(string(json_data), nil))
	s.mock.On("Exists", mock.Anything, mock.Anything).Return(redis.NewIntResult(1, nil))
	res, err := s.repo.Fetch(context.Background(), 0)
	s.Assert().Equal(err, nil)
	s.Assert().Equal(res, IDs)
}

func (s *CharacterReadRepositoryTestSuite) TestSuccessFetch() {
	IDs := []int{1, 2, 3}
	json_data, err := json.Marshal(IDs)
	if err != nil {
		s.T().Fatalf("Error: '%s'", err)
	}

	s.mock.On("Get", mock.Anything, "marvel-characters-page-3").Return(redis.NewStringResult(string(json_data), nil))
	s.mock.On("Exists", mock.Anything, mock.Anything).Return(redis.NewIntResult(1, nil))
	res, err := s.repo.Fetch(context.Background(), 3)
	s.Assert().Equal(err, nil)
	s.Assert().Equal(res, IDs)
}

func (s *CharacterReadRepositoryTestSuite) TestNilGetByID() {
	s.mock.On("Get", mock.Anything, "marvel-character-id-1").Return(redis.NewStringResult("", nil))
	s.mock.On("Exists", mock.Anything, mock.Anything).Return(redis.NewIntResult(1, nil))

	_, err := s.repo.GetByID(context.Background(), 1)
	s.Assert().Equal(err, domain.ErrNotFound)
}

func (s *CharacterReadRepositoryTestSuite) TestFailedJSONGetByID() {
	s.mock.On("Exists", mock.Anything, mock.Anything).Return(redis.NewIntResult(1, nil))
	s.mock.On("Get", mock.Anything, "marvel-character-id-2").Return(redis.NewStringResult("val", nil))

	_, err := s.repo.GetByID(context.Background(), 2)
	s.Assert().Equal(err, domain.ErrInternalServerError)
}

func (s *CharacterReadRepositoryTestSuite) TestEmptyKeyGetByID() {
	s.mock.On("Exists", mock.Anything, mock.Anything).Return(redis.NewIntResult(0, nil))

	_, err := s.repo.GetByID(context.Background(), 2)
	s.Assert().Equal(err, domain.ErrCacheKeyEmpty)
}

func (s *CharacterReadRepositoryTestSuite) TestErrorEmptyKeyGetByID() {
	s.mock.On("Exists", mock.Anything, mock.Anything).Return(redis.NewIntResult(1, errors.New("err")))

	_, err := s.repo.GetByID(context.Background(), 2)
	s.Assert().Equal(err, domain.ErrInternalServerError)
}

func (s *CharacterReadRepositoryTestSuite) TestSuccessGetByID() {
	record = domain.Character{
		ID:          1,
		Name:        "lorem",
		Description: "lorem",
		FetchedAt:   time.Now(),
	}
	json_data, err := json.Marshal(record)
	if err != nil {
		s.T().Fatalf("Error: '%s'", err)
	}

	s.mock.On("Exists", mock.Anything, mock.Anything).Return(redis.NewIntResult(1, nil))
	s.mock.On("Get", mock.Anything, "marvel-character-id-3").Return(redis.NewStringResult(string(json_data), nil))

	res, err := s.repo.GetByID(context.Background(), 3)
	s.Assert().Equal(err, nil)
	s.Assert().Equal(res.ID, record.ID)
	s.Assert().Equal(res.Name, record.Name)
	s.Assert().Equal(res.FetchedAt.Format(time.RFC3339Nano), record.FetchedAt.Format(time.RFC3339Nano))
}
