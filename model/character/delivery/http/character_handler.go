package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/hezbymuhammad/golang-marvel-demo/domain"
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
	case domain.ErrCacheKeyEmpty:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

type CharacterHandler struct {
	Usecase domain.CharacterUsecase
}

func NewCharacterHandler(e *echo.Echo, u domain.CharacterUsecase) *CharacterHandler {
	handler := &CharacterHandler{
		Usecase: u,
	}
	e.GET("/characters", handler.Fetch)
	e.GET("/characters/", handler.Fetch)
	e.GET("/characters/:id", handler.GetByID)

	return handler
}

func (h *CharacterHandler) Fetch(c echo.Context) error {
	pageRaw := c.QueryParam("page")
	if pageRaw == "" {
		pageRaw = "0"
	}

	page, err := strconv.Atoi(pageRaw)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Bad request param"})
	}

	ctx := c.Request().Context()

	IDs, err := h.Usecase.Fetch(ctx, page)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(getStatusCode(err), IDs)
}

func (h *CharacterHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Bad request param"})
	}

	ctx := c.Request().Context()

	character, err := h.Usecase.GetByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(getStatusCode(err), character)
}
