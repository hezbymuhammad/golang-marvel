package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	redismock "github.com/elliotchance/redismock/v8"
	redis "github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	gock "gopkg.in/h2non/gock.v1"

	"github.com/hezbymuhammad/golang-marvel/domain"
	"github.com/hezbymuhammad/golang-marvel/model/character/repository"
)

type CharacterWriteRepositoryTestSuite struct {
	suite.Suite
	redisMock *redismock.ClientMock
	repo      domain.CharacterWriteRepository
}

func TestCharacterWriteRepository(t *testing.T) {
	suite.Run(t, new(CharacterWriteRepositoryTestSuite))
}

func (s *CharacterWriteRepositoryTestSuite) SetupTest() {
	defer gock.Off()

	api, pubK, privK := "http://foo.com", "asd", "asd"
	timeout := 2 * time.Second
	cacheExpiration := 10 * time.Second

	mr, err := miniredis.Run()
	if err != nil {
		s.T().Fatalf("Error: '%s'", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	s.redisMock = redismock.NewNiceMock(client)
	s.repo = repository.NewCharacterWriteRepository(api, pubK, privK, s.redisMock, timeout, cacheExpiration)
}

func (s *CharacterWriteRepositoryTestSuite) TestSuccessStoreByPage() {
	gock.New("http://foo.com").Get("/v1/public/characters/").Reply(200).BodyString("{\"data\": { \"results\": [{\"id\": 1011334, \"name\": \"lorem\", \"description\": \"asd\"}] }}")
	s.redisMock.On("Set", mock.Anything, "marvel-characters-page-1", "[1011334]", mock.Anything).Return(redis.NewStatusResult("", nil))
	s.redisMock.On("Set", mock.Anything, "marvel-character-id-1011334", mock.Anything, mock.Anything).Return(redis.NewStatusResult("", nil))

	err := s.repo.StoreByPage(context.Background(), 1)
	s.Assert().Equal(err, nil)
}

func (s *CharacterWriteRepositoryTestSuite) TestNilDataStoreByPage() {
	gock.New("http://foo.com").Get("/v1/public/characters/").Reply(200).BodyString("{\"data\": {  }}")

	err := s.repo.StoreByPage(context.Background(), 1)
	s.Assert().Equal(err, domain.ErrNotFound)
}

func (s *CharacterWriteRepositoryTestSuite) TestFailedJSONStoreByPage() {
	gock.New("http://foo.com").Get("/v1/public/characters/").Reply(200).BodyString("val")

	err := s.repo.StoreByPage(context.Background(), 1)
	s.Assert().Equal(err, domain.ErrInternalServerError)
}

func (s *CharacterWriteRepositoryTestSuite) TestFailedRedisStoreByPage() {
	gock.New("http://foo.com").Get("/v1/public/characters/").Reply(200).BodyString("{\"data\": { \"results\": [{\"id\": 1011334, \"name\": \"lorem\", \"description\": \"asd\"}] }}")
	s.redisMock.On("Set", mock.Anything, "marvel-characters-page-1", "[1011334]", mock.Anything).Return(redis.NewStatusResult("", errors.New("error")))
	s.redisMock.On("Set", mock.Anything, "marvel-character-id-1011334", mock.Anything, mock.Anything).Return(redis.NewStatusResult("", nil))

	err := s.repo.StoreByPage(context.Background(), 1)
	s.Assert().Equal(err, domain.ErrInternalServerError)
}

func (s *CharacterWriteRepositoryTestSuite) TestSuccessStoreByID() {
	gock.New("http://foo.com").Get("/v1/public/characters/1").Reply(200).BodyString("{\"data\": { \"results\": [{\"id\": 1011334, \"name\": \"lorem\", \"description\": \"asd\"}] }}")
	s.redisMock.On("Set", mock.Anything, "marvel-character-id-1011334", mock.Anything, mock.Anything).Return(redis.NewStatusResult("", nil))

	err := s.repo.StoreByID(context.Background(), 1)
	s.Assert().Equal(err, nil)
}

func (s *CharacterWriteRepositoryTestSuite) TestNilDataStoreByID() {
	gock.New("http://foo.com").Get("/v1/public/characters/1").Reply(200).BodyString("{\"data\": {  }}")

	err := s.repo.StoreByID(context.Background(), 1)
	s.Assert().Equal(err, domain.ErrNotFound)
}

func (s *CharacterWriteRepositoryTestSuite) TestFailedJSONStoreByID() {
	gock.New("http://foo.com").Get("/v1/public/characters/1").Reply(200).BodyString("val")

	err := s.repo.StoreByID(context.Background(), 1)
	s.Assert().Equal(err, domain.ErrInternalServerError)
}

func (s *CharacterWriteRepositoryTestSuite) TestFailedRedisStoreByID() {
	gock.New("http://foo.com").Get("/v1/public/characters/1").Reply(200).BodyString("{\"data\": { \"results\": [{\"id\": 1011334, \"name\": \"lorem\", \"description\": \"asd\"}] }}")
	s.redisMock.On("Set", mock.Anything, "marvel-character-id-1011334", mock.Anything, mock.Anything).Return(redis.NewStatusResult("", errors.New("error")))

	err := s.repo.StoreByID(context.Background(), 1)
	s.Assert().Equal(err, domain.ErrInternalServerError)
}
