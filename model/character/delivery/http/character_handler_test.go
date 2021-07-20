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

	"hezbymuhammad/golang-marvel-demo/domain"
	"hezbymuhammad/golang-marvel-demo/domain/mocks"
	characterHttp "hezbymuhammad/golang-marvel-demo/model/character/delivery/http"
)

type CharacterHandlerTestSuite struct {
	suite.Suite
        handler       characterHttp.CharacterHandler
	usecase       *mocks.CharacterUsecase
}

func TestCharacterHandler(t *testing.T) {
	suite.Run(t, new(CharacterHandlerTestSuite))
}

func (s *CharacterHandlerTestSuite) SetupTest() {
	s.usecase = new(mocks.CharacterUsecase)
        s.handler = characterHttp.CharacterHandler{
                Usecase: s.usecase,
        }
}

func (s *CharacterHandlerTestSuite) TestSuccessFetch() {
        e := echo.New()
        rec := httptest.NewRecorder()
        req, err := http.NewRequest(echo.GET, "/characters?page=1", strings.NewReader(""))
        ctx := e.NewContext(req, rec)

        arr := []int{1, 2, 3}
        s.usecase.On("Fetch", mock.Anything, 1).Return(arr, nil).Once()

        err = s.handler.Fetch(ctx)
        s.Assert().Equal(err, nil)
	s.Assert().Equal(http.StatusOK, rec.Code)
}

func (s *CharacterHandlerTestSuite) TestFailedFetch() {
        e := echo.New()
        rec := httptest.NewRecorder()
        req, err := http.NewRequest(echo.GET, "/characters?page=1", strings.NewReader(""))
        ctx := e.NewContext(req, rec)

        arr := []int{1, 2, 3}
        s.usecase.On("Fetch", mock.Anything, 1).Return(arr, errors.New("SomeError")).Once()

        err = s.handler.Fetch(ctx)
        s.Assert().Equal(err, nil)
	s.Assert().Equal(http.StatusInternalServerError, rec.Code)
	s.Assert().Equal("{\"message\":\"SomeError\"}\n", rec.Body.String())
}

func (s *CharacterHandlerTestSuite) TestSuccessGetByID() {
        record := domain.Character{
                ID: 1,
                Name: "Lorem",
                Description: "Lorem",
                FetchedAt: time.Now(),
        }
        json_data, err := json.Marshal(record)
        id := int(record.ID)

        e := echo.New()
        rec := httptest.NewRecorder()
        req, err := http.NewRequest(echo.GET, "/characters/" + strconv.Itoa(id), strings.NewReader(""))
        ctx := e.NewContext(req, rec)
        ctx.SetPath("characters/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(strconv.Itoa(id))

        s.usecase.On("GetByID", mock.Anything, id).Return(record, nil).Once()

        err = s.handler.GetByID(ctx)
        s.Assert().Equal(err, nil)
	s.Assert().Equal(http.StatusOK, rec.Code)
	s.Assert().Equal(string(json_data) + "\n", rec.Body.String())
}

func (s *CharacterHandlerTestSuite) TestFailedGetByID() {
        record := domain.Character{
                ID: 1,
                Name: "Lorem",
                Description: "Lorem",
                FetchedAt: time.Now(),
        }
        id := int(record.ID)

        e := echo.New()
        rec := httptest.NewRecorder()
        req, err := http.NewRequest(echo.GET, "/characters/" + strconv.Itoa(id), strings.NewReader(""))
        ctx := e.NewContext(req, rec)
        ctx.SetPath("characters/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(strconv.Itoa(id))

        s.usecase.On("GetByID", mock.Anything, id).Return(record, errors.New("SomeError")).Once()

        err = s.handler.GetByID(ctx)
        s.Assert().Equal(err, nil)
	s.Assert().Equal(http.StatusInternalServerError, rec.Code)
	s.Assert().Equal("{\"message\":\"SomeError\"}\n", rec.Body.String())
}
