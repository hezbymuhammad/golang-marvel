package http_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/hezbymuhammad/golang-marvel-demo/domain"
	"github.com/hezbymuhammad/golang-marvel-demo/domain/mocks"
	characterHttp "github.com/hezbymuhammad/golang-marvel-demo/model/character/delivery/http"
)

type CharacterHandlerTestSuite struct {
	suite.Suite
	handler *characterHttp.CharacterHandler
	usecase *mocks.CharacterUsecase
}

func TestCharacterHandler(t *testing.T) {
	suite.Run(t, new(CharacterHandlerTestSuite))
}

func (s *CharacterHandlerTestSuite) SetupTest() {
	s.usecase = new(mocks.CharacterUsecase)
	s.handler = characterHttp.NewCharacterHandler(echo.New(), s.usecase)

}

func (s *CharacterHandlerTestSuite) TestSuccessFetch() {
	e := echo.New()
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(echo.GET, "/characters?page=2", strings.NewReader(""))
	ctx := e.NewContext(req, rec)

	arr := []int{1, 2, 3}
	s.usecase.On("Fetch", mock.Anything, 2).Return(arr, nil)

	err = s.handler.Fetch(ctx)
	s.Assert().Equal(err, nil)
	s.Assert().Equal(http.StatusOK, rec.Code)
}

func (s *CharacterHandlerTestSuite) TestPageNilFetch() {
	e := echo.New()
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(echo.GET, "/characters/", strings.NewReader(""))
	ctx := e.NewContext(req, rec)

	arr := []int{1, 2, 3}
	s.usecase.On("Fetch", mock.Anything, 0).Return(arr, nil)

	err = s.handler.Fetch(ctx)
	s.Assert().Equal(err, nil)
	s.Assert().Equal(http.StatusOK, rec.Code)
}

func (s *CharacterHandlerTestSuite) TestWrongPageFetch() {
	e := echo.New()
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(echo.GET, "/characters?page=aaaa", strings.NewReader(""))
	ctx := e.NewContext(req, rec)

	err = s.handler.Fetch(ctx)
	s.Assert().Equal(err, nil)
	s.Assert().Equal(http.StatusBadRequest, rec.Code)
	s.Assert().Equal("{\"message\":\"Bad request param\"}\n", rec.Body.String())
}

func (s *CharacterHandlerTestSuite) TestFailedFetch() {
	e := echo.New()
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(echo.GET, "/characters?page=1", strings.NewReader(""))
	ctx := e.NewContext(req, rec)

	arr := []int{1, 2, 3}
	s.usecase.On("Fetch", mock.Anything, 1).Return(arr, errors.New("SomeError"))

	err = s.handler.Fetch(ctx)
	s.Assert().Equal(err, nil)
	s.Assert().Equal(http.StatusInternalServerError, rec.Code)
	s.Assert().Equal("{\"message\":\"SomeError\"}\n", rec.Body.String())
}

func (s *CharacterHandlerTestSuite) TestSuccessGetByID() {
	e := echo.New()
	rec := httptest.NewRecorder()
	record := domain.Character{
		ID:          1,
		Name:        "Lorem",
		Description: "Lorem",
		FetchedAt:   time.Now(),
	}
	json_data, err := json.Marshal(record)
	id := int(record.ID)

	req, err := http.NewRequest(echo.GET, "/characters/"+strconv.Itoa(id), strings.NewReader(""))
	ctx := e.NewContext(req, rec)
	ctx.SetPath("characters/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(strconv.Itoa(id))

	s.usecase.On("GetByID", mock.Anything, id).Return(record, nil)

	err = s.handler.GetByID(ctx)
	s.Assert().Equal(err, nil)
	s.Assert().Equal(http.StatusOK, rec.Code)
	s.Assert().Equal(string(json_data)+"\n", rec.Body.String())
}

func (s *CharacterHandlerTestSuite) TestFailedGetByID() {
	e := echo.New()
	rec := httptest.NewRecorder()
	record := domain.Character{
		ID:          2,
		Name:        "Lorem",
		Description: "Lorem",
		FetchedAt:   time.Now(),
	}
	id := int(record.ID)

	req, err := http.NewRequest(echo.GET, "/characters/"+strconv.Itoa(id), strings.NewReader(""))
	ctx := e.NewContext(req, rec)
	ctx.SetPath("characters/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(strconv.Itoa(id))

	s.usecase.On("GetByID", mock.Anything, id).Return(record, errors.New("SomeError"))

	err = s.handler.GetByID(ctx)
	s.Assert().Equal(err, nil)
	s.Assert().Equal(http.StatusInternalServerError, rec.Code)
	s.Assert().Equal("{\"message\":\"SomeError\"}\n", rec.Body.String())
}

func (s *CharacterHandlerTestSuite) TestNotFoundGetByID() {
	e := echo.New()
	rec := httptest.NewRecorder()
	record := domain.Character{
		ID:          3,
		Name:        "Lorem",
		Description: "Lorem",
		FetchedAt:   time.Now(),
	}
	id := int(record.ID)

	req, err := http.NewRequest(echo.GET, "/characters/"+strconv.Itoa(id), strings.NewReader(""))
	ctx := e.NewContext(req, rec)
	ctx.SetPath("characters/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(strconv.Itoa(id))

	s.usecase.On("GetByID", mock.Anything, id).Return(record, domain.ErrNotFound)

	err = s.handler.GetByID(ctx)
	s.Assert().Equal(err, nil)
	s.Assert().Equal(http.StatusNotFound, rec.Code)
	s.Assert().Equal("{\"message\":\"Resource not found\"}\n", rec.Body.String())
}

func (s *CharacterHandlerTestSuite) TestWrongIDGetByID() {
	e := echo.New()
	rec := httptest.NewRecorder()

	req, err := http.NewRequest(echo.GET, "/characters/", strings.NewReader(""))
	ctx := e.NewContext(req, rec)

	err = s.handler.GetByID(ctx)
	s.Assert().Equal(err, nil)
	s.Assert().Equal(http.StatusBadRequest, rec.Code)
	s.Assert().Equal("{\"message\":\"Bad request param\"}\n", rec.Body.String())
}
