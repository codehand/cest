package test

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HealthCheck4(c echo.Context) error {
	// To do something
	return c.JSON(http.StatusOK, 1)
}

// golang
