package usecase_test

import (
        "context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
        "github.com/stretchr/testify/suite"

	"github.com/hezbymuhammad/golang-marvel/domain"
	"github.com/hezbymuhammad/golang-marvel/domain/mocks"
	"github.com/hezbymuhammad/golang-marvel/model/character/usecase"
)

type CharacterUsecaseTestSuite struct {
	suite.Suite
	usecase       domain.CharacterUsecase
	readRepo      *mocks.CharacterReadRepository
	writeRepo     *mocks.CharacterWriteRepository
}

func TestCharacterUsecase(t *testing.T) {
	suite.Run(t, new(CharacterUsecaseTestSuite))
}

func (s *CharacterUsecaseTestSuite) SetupTest() {
	s.readRepo = new(mocks.CharacterReadRepository)
	s.writeRepo = new(mocks.CharacterWriteRepository)
        s.usecase = usecase.NewCharacterUsecase(s.readRepo, s.writeRepo, time.Second*2)
}

func (s *CharacterUsecaseTestSuite) TestSuccessFetch() {
        arr := []int{1, 2, 3}
        s.readRepo.On("Fetch", mock.Anything, 1).Return(arr, nil).Once()
        s.writeRepo.On("StoreByPage", mock.Anything, 1).Return(nil).Once()

        res, err := s.usecase.Fetch(context.Background(), 1)
        s.Assert().Equal(res, arr)
        s.Assert().Equal(err, nil)
}

func (s *CharacterUsecaseTestSuite) TestFailedFetch() {
        arr := []int{1, 2, 3}
        dummy_err := errors.New("SomeError")
        s.readRepo.On("Fetch", mock.Anything, 1).Return(arr, dummy_err).Once()
        s.writeRepo.On("StoreByPage", mock.Anything, 1).Return(nil).Once()

        res, err := s.usecase.Fetch(context.Background(), 1)
        s.Assert().Nil(res)
        s.Assert().Equal(err, dummy_err)
}

func (s *CharacterUsecaseTestSuite) TestFailedStoreByPage() {
        arr := []int{1, 2, 3}
        s.readRepo.On("Fetch", mock.Anything, 1).Return(arr, nil).Once()
        s.writeRepo.On("StoreByPage", mock.Anything, 1).Return(errors.New("SomeError")).Once()

        res, err := s.usecase.Fetch(context.Background(), 1)
        s.Assert().Equal(res, arr)
        s.Assert().Equal(err, nil)
}

func (s *CharacterUsecaseTestSuite) TestSuccessGetByID() {
        record := domain.Character{
                ID: 1,
                Name: "Lorem",
                Description: "Lorem",
                FetchedAt: time.Now(),
        }
        s.readRepo.On("GetByID", mock.Anything, 1).Return(record, nil).Once()
        s.writeRepo.On("StoreByID", mock.Anything, 1).Return(nil).Once()

        res, err := s.usecase.GetByID(context.Background(), 1)
        s.Assert().Equal(res, record)
        s.Assert().Equal(err, nil)
}

func (s *CharacterUsecaseTestSuite) TestFailedGetByID() {
        record := domain.Character{
                ID: 1,
                Name: "Lorem",
                Description: "Lorem",
                FetchedAt: time.Now(),
        }
        dummy_err := errors.New("SomeError")
        s.readRepo.On("GetByID", mock.Anything, 1).Return(record, dummy_err).Once()
        s.writeRepo.On("StoreByID", mock.Anything, 1).Return(nil).Once()

        res, err := s.usecase.GetByID(context.Background(), 1)
        s.Assert().Equal(res, domain.Character{})
        s.Assert().Equal(err, dummy_err)
}

func (s *CharacterUsecaseTestSuite) TestFailedStoreByID() {
        record := domain.Character{
                ID: 1,
                Name: "Lorem",
                Description: "Lorem",
                FetchedAt: time.Now(),
        }
        s.readRepo.On("GetByID", mock.Anything, 1).Return(record, nil).Once()
        s.writeRepo.On("StoreByID", mock.Anything, 1).Return(errors.New("SomeError")).Once()

        res, err := s.usecase.GetByID(context.Background(), 1)
        s.Assert().Equal(res, record)
        s.Assert().Equal(err, nil)
}
