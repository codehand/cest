package test

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ABC is func1
// fn testing
// func ABC(in int) error {
// 	return nil
// }

// // XYZ is abc
// func XYZ(c echo.Context) error {
// 	return nil
// }

func HealthBeforeCheck(c echo.Context) error {
	// To do something
	return c.JSON(http.StatusOK, true)
}

func HealthCheck(c echo.Context) error {
	// To do something
	return c.JSON(http.StatusOK, "true")
}

func HealthCheck2(c echo.Context) error {
	// To do something
	return c.JSON(http.StatusOK, 1)
}

// golang
