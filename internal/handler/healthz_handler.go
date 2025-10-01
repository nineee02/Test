package handler

import "github.com/labstack/echo/v4"

type HealthHandler interface {
	Check(c echo.Context) error
}

type healthHandler struct{}

func NewHealthHandler() HealthHandler {
	return &healthHandler{}
}

func (h *healthHandler) Check(c echo.Context) error {
	return c.JSON(200, echo.Map{
		"status":  "ok",
		"message": "API user-api is running and healthy",
	})
}
