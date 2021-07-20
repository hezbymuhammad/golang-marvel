package http

import (
	"net/http"
	"strconv"
        "fmt"

	"github.com/labstack/echo"

	"hezbymuhammad/golang-marvel-demo/domain"
)

type ResponseError struct {
	Message string `json:"message"`
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

        fmt.Errorf("HTTP Error: %v", err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

type CharacterHandler struct {
	Usecase domain.CharacterUsecase
}

func NewCharacterHandler(e *echo.Echo, u domain.CharacterUsecase) {
	handler := &CharacterHandler{
		Usecase: u,
	}
	e.GET("/characters", handler.Fetch)
	e.GET("/characters/:id", handler.GetByID)
}

func (h *CharacterHandler) Fetch(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	ctx := c.Request().Context()

	IDs, err := h.Usecase.Fetch(ctx, page)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, IDs)
}

func (h *CharacterHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	ctx := c.Request().Context()

	character, err := h.Usecase.GetByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, character)
}
